package ui

import "github.com/charmbracelet/lipgloss"

var (
	primaryColor = lipgloss.Color("#36968e")
	subtleColor  = lipgloss.Color("241")
	alertColor   = lipgloss.Color("196")

	TimeStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	PausedStyle = lipgloss.NewStyle().
			Foreground(subtleColor).
			Bold(true)

	AlarmStyle = lipgloss.NewStyle().
			Foreground(alertColor).
			Bold(true).
			Blink(true)

	BarFilledStyle = lipgloss.NewStyle().Foreground(primaryColor)
	BarEmptyStyle  = lipgloss.NewStyle().Foreground(subtleColor)
)
