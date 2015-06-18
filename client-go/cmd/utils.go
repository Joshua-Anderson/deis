package cmd

import (
	"fmt"
	"strings"
	"time"
)

func progress() chan bool {
	frames := []string{"...", "o..", ".o.", "..o"}
	backspaces := strings.Repeat("\b", 3)
	tick := time.Tick(400 * time.Millisecond)
	quit := make(chan bool)
	go func() {
		for {
			for _, frame := range frames {
				fmt.Print(frame)
				select {
				case <-quit:
					fmt.Print(backspaces)
					close(quit)
					return
				case <-tick:
					fmt.Print(backspaces)
				}
			}
		}
	}()
	return quit
}

// Chose a ansi color based on converting a string to a int.
func chooseColor(input string) string {
	var sum uint8

	for _, char := range []byte(input) {
		sum += uint8(char)
	}

	// Seven possible terminal colors
	color := (sum % 7) + 1

	if color == 7 {
		color = 9
	}

	return fmt.Sprintf("\033[3%dm", color)
}

func printColor(line string, ANSIColor string) {
	fmt.Printf("%s%s\033[39m\n", ANSIColor, line)
}
