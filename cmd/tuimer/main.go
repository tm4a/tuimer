package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tm4a/tuimer/internal/app"
)

func printHelp() {
	fmt.Println("Tuimer - Minimal TUI Timer")
	fmt.Println("\nUSAGE:")
	fmt.Println("  tuimer              Start interactive mode")
	fmt.Println("  tuimer <duration>   Quick start (e.g. '10m', '1h30m', '45s')")
	fmt.Println("  tuimer --help       Show this message")
	fmt.Println("\nCONTROLS:")
	fmt.Println("  0-9         Type time (HH:MM:SS)")
	fmt.Println("  Enter       Start timer / Stop alarm (Reset)")
	fmt.Println("  Space       Pause / Resume")
	fmt.Println("  q           Reset to input (if running) / Quit (if at input)")
	fmt.Println("  Ctrl+C      Force Quit")
	fmt.Println("\nHIDDEN SHORTCUTS (Input Mode):")
	fmt.Println("  p           Pomodoro (25m)")
	fmt.Println("  s           Short Break (5m)")
	fmt.Println("  l           Long Break (15m)")
	fmt.Println("\nCONFIGURATION:")
	fmt.Println("  Place alarm.mp3 (or .wav/.ogg) in ~/.config/tuimer/")
}

func main() {
	startSeconds := 0

	if len(os.Args) > 1 {
		arg := os.Args[1]
		if arg == "-h" || arg == "--help" || arg == "help" {
			printHelp()
			os.Exit(0)
		}

		input := strings.Join(os.Args[1:], "")
		d, err := time.ParseDuration(input)
		if err != nil {
			fmt.Println("Error: Invalid time format.")
			fmt.Println("Run 'tuimer --help' for details.")
			os.Exit(1)
		}
		startSeconds = int(d.Seconds())
	}

	p := tea.NewProgram(app.NewModel(startSeconds), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}
}
