package timer

import (
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type TickMsg time.Time

func Tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func FormatBuffer(buf string) string {
	padded := buf
	for len(padded) < 6 {
		padded = "0" + padded
	}
	return padded[0:2] + ":" + padded[2:4] + ":" + padded[4:6]
}

func ParseAndValidateTime(buf string) int {
	padded := buf
	for len(padded) < 6 {
		padded = "0" + padded
	}
	h, _ := strconv.Atoi(padded[0:2])
	m, _ := strconv.Atoi(padded[2:4])
	s, _ := strconv.Atoi(padded[4:6])
	if m > 59 {
		m = 59
	}
	if s > 59 {
		s = 59
	}
	return (h * 3600) + (m * 60) + s
}
