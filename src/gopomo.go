package main

import (
	"fmt"
	"os"
	"time"
)

const LOOPS_AFTER_BREAK = 3

func usage_error() {
  fmt.Printf("Usage: %s <work-time(m)> <break-time(m)> <optional confirm-flag 1/0>\n", os.Args[0])
  os.Exit(0)
}

func confirmPrompt() {
  fmt.Println("Confirming next stage, press Enter to continue...")
  fmt.Scanln()
}

func main() {
	var argc int = len(os.Args)
  var confirmFlag bool = false

  if argc == 4 {
    if os.Args[3] == "1" {
      confirmFlag = true
      fmt.Println("Asking to confirm.")
    } else if os.Args[3] == "0" {
      confirmFlag = false
      fmt.Println("Confirmation disabled.")
    } else {
      usage_error()
    }
  } else if argc < 3 || argc > 4 {
    usage_error()
  }

	// Parse arguments
	workStr := os.Args[1]
	workTime, err := time.ParseDuration(workStr)
	if err != nil {
		fmt.Printf("Invalid work time format: %v\n", err)
		os.Exit(1)
	}

	breakStr := os.Args[2]
	breakTime, err := time.ParseDuration(breakStr)
	if err != nil {
		fmt.Printf("Invalid break time format: %v\n", err)
		os.Exit(1)
	}

	var isBreak bool = false
	var longBreakTimer int = 0

	for {
		if !isBreak { // Work section
      if confirmFlag == true {
        confirmPrompt()
      }

			fmt.Println("Starting Work.")
			ticker := time.NewTicker(1 * time.Second)
			startTime := time.Now()

			go func() {
				for range ticker.C {
					elapsed := time.Since(startTime)
					remaining := workTime - elapsed
					if remaining <= 0 {
						ticker.Stop()
						return
					}
					fmt.Printf("Remaining work time: %v\n", remaining.Round(time.Second))
				}
			}()

			time.Sleep(workTime)
			fmt.Println("Work is Over.")
			isBreak = true
		} else if isBreak && longBreakTimer != LOOPS_AFTER_BREAK { // Break section
			if confirmFlag == true {
        confirmPrompt()
      }

      fmt.Println("Starting Break.")
			ticker := time.NewTicker(1 * time.Second)
			startTime := time.Now()

			go func() {
				for range ticker.C {
					elapsed := time.Since(startTime)
					remaining := breakTime - elapsed
					if remaining <= 0 {
						ticker.Stop()
						return
					}
					fmt.Printf("Remaining break time: %v\n", remaining.Round(time.Second))
				}
			}()

			time.Sleep(breakTime)
			fmt.Println("Break is Over.")
			isBreak = false
			longBreakTimer++
		} else if isBreak && longBreakTimer == LOOPS_AFTER_BREAK { // Long break section
			if confirmFlag == true {
        confirmPrompt()
      }

      fmt.Println("Starting Long Break.")
			ticker := time.NewTicker(1 * time.Second)
			startTime := time.Now()

			go func() {
				for range ticker.C {
					elapsed := time.Since(startTime)
					remaining := (breakTime * 2) - elapsed
					if remaining <= 0 {
						ticker.Stop()
						return
					}
					fmt.Printf("Remaining long break time: %v\n", remaining.Round(time.Second))
				}
			}()

			time.Sleep(breakTime * 2)
			ticker.Stop()
			fmt.Println("Long Break is Over.")
			isBreak = false
			longBreakTimer = 0
		}

	}
}
