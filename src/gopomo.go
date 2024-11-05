package main

import (
	"bufio"
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
	flag.IntVar(&cfg.breakLoops, "loops", 3, "Defined as number of work sessions before long break")
	flag.IntVar(&cfg.breakLoops, "l", 3, "Loops of work before long break, shorthand")
	flag.BoolVar(&cfg.confirmBreak, "confirm", false, "Confirm before starting next phase")
	flag.BoolVar(&cfg.confirmBreak, "c", false, "Confirm before starting next phase, shorthand")

	// Custom usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage for %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "A simple Pomodoro app.\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample Usage:\n")
		fmt.Fprintf(os.Stderr, "%s -word 25 -break 5 -loops 6 -confirm\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s -w 45 -b 7.5 -lb 22 -l 5 -c\n", os.Args[0])
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
			cfg.longBreakTime = time.Duration((longBreakMinutes) * float64(time.Minute))
		}
	})

	return cfg, nil
}

func runPomodoro(cfg config) {
	isBreak := false
	longBreakCounter := 0

	pauseChan := make(chan bool)
	go listenForPause(pauseChan)

	for {
		var duration time.Duration
		var phase string

		if !isBreak {
			duration = cfg.workTime
			phase = "Work"
			longBreakCounter++
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
			<-pauseChan
		}

		fmt.Printf("\nStarting %s phase.\n", phase)
		countdown(duration, pauseChan)
		fmt.Printf("%s complete.\a\n", phase)

		isBreak = !isBreak
	}
}

func countdown(duration time.Duration, pauseChan <-chan bool) {
	remaining := duration
	paused := false
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	fmt.Printf("Press Return to pause the timer.\n")

	lastTick := time.Now()
	for remaining > 0 {
		select {
		case <-ticker.C:
			if !paused {
				now := time.Now()
				elapsed := now.Sub(lastTick)
				remaining -= elapsed
				lastTick = now
				fmt.Printf("\rTime remaining: %v  ", remaining.Round(time.Second))
			}
		case pause := <-pauseChan:
			if pause {
				if !paused {
					paused = true
					fmt.Println("Timer paused. Press Return to resume.")
				} else {
					paused = false
					lastTick = time.Now()
					fmt.Println("Timer resumed.")
				}
			}
		}
	}
	fmt.Println()
}

func listenForPause(pauseChan chan<- bool) {
	reader := bufio.NewReader(os.Stdin)
	for {
		_, _ = reader.ReadString('\n')
		pauseChan <- true
	}
}
