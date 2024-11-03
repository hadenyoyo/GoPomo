package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

type config struct {
	workTime      time.Duration
	breakTime     time.Duration
	longBreakTime time.Duration
	breakLoops    int
	confirmBreak  bool
}

func main() {
	cfg, err := parseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Starting GoPomo with work time: %v, short break time %v, and long break time %v\n", cfg.workTime, cfg.breakTime, cfg.longBreakTime)

	runPomodoro(cfg)
}

func parseFlags() (config, error) {
	var cfg config
	var workMinutes, breakMinutes, longBreakMinutes float64

	flag.Float64Var(&workMinutes, "work", 25, "Work time in minutes")
	flag.Float64Var(&workMinutes, "w", 25, "Work time in minutes, shorthand")
	flag.Float64Var(&breakMinutes, "break", 5, "Break time in minutes")
	flag.Float64Var(&breakMinutes, "b", 5, "Break time in minutes, shorthand")
	flag.Float64Var(&longBreakMinutes, "longbreak", 0, "Long break time in minutes")
	flag.Float64Var(&longBreakMinutes, "lb", 0, "Long break time in minutes, shorthand")
	flag.IntVar(&cfg.breakLoops, "loops", 3, "Loops of work/break before long break")
	flag.IntVar(&cfg.breakLoops, "l", 3, "Loops of work/break before long break, shorthand")
	flag.BoolVar(&cfg.confirmBreak, "confirm", false, "Confirm before starting next phase")
	flag.BoolVar(&cfg.confirmBreak, "c", false, "Confirm before starting next phase, shorthand")

	// Custom usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage for %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "A simple Pomodoro app.\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample Usage:\n")
		fmt.Fprintf(os.Stderr, "%s -word 25 -break 5 -confirm\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s -w 45 -b 7 -c\n", os.Args[0])
	}

	flag.Parse()

	// Convert to time format
	cfg.workTime = time.Duration((workMinutes) * float64(time.Minute))
	cfg.breakTime = time.Duration((breakMinutes) * float64(time.Minute))

	// longBreakTime default value
	cfg.longBreakTime = cfg.breakTime * 2

	// Update longBreakTime if explicitly set
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "longbreak" || f.Name == "lb" {
			cfg.longBreakTime = time.Duration(longBreakMinutes) * time.Minute
		}
	})

	return cfg, nil
}

func runPomodoro(cfg config) {
	isBreak := false
	longBreakCounter := 0

	for {
		var duration time.Duration
		var phase string

		if !isBreak {
			duration = cfg.workTime
			phase = "Work"
		} else if longBreakCounter < cfg.breakLoops {
			duration = cfg.breakTime
			phase = "Break"
		} else {
			duration = cfg.longBreakTime
			phase = "Long Break"
			longBreakCounter = 0
		}

		if cfg.confirmBreak {
			fmt.Println("Confirming next stage, press Enter to continue...")
			fmt.Scanln()
		}

		fmt.Printf("Starting %s phase.\n", phase)
		countdown(duration)
		fmt.Printf("%s complete.\n", phase)

		if isBreak {
			longBreakCounter++
		}
		isBreak = !isBreak
	}
}

func countdown(duration time.Duration) {
	start := time.Now()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for remaining := duration; remaining > 0; remaining = duration - time.Since(start) {
		fmt.Printf("Time remaining: %v\n", remaining.Round(time.Second))
		<-ticker.C
	}
}
