package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/tm4a/tuimer/internal/timer"
	"github.com/tm4a/tuimer/internal/ui"
)

func (m Model) View() string {
	if m.Quitting {
		return ""
	}

	var timeStr string
	if m.InputMode {
		timeStr = timer.FormatBuffer(m.InputBuffer)
		timeStr = ui.TimeStyle.Render(timeStr)
	} else {
		h := m.TimeLeft / 3600
		min := (m.TimeLeft % 3600) / 60
		sec := m.TimeLeft % 60

		rawTime := fmt.Sprintf("%02d:%02d:%02d", h, min, sec)

		if m.Finished {
			timeStr = ui.AlarmStyle.Render(rawTime)
		} else if m.Paused {
			timeStr = ui.PausedStyle.Render(rawTime + " [PAUSED]")
		} else {
			timeStr = ui.TimeStyle.Render(rawTime)
		}
	}

	var bottomStr string

	if !m.InputMode && !m.Finished {
		barWidth := int(float64(m.Width) * 0.6)
		if barWidth < 10 {
			barWidth = 10
		}
		if barWidth > 100 {
			barWidth = 100
		}

		filledWidth := int(float64(barWidth) * m.Percent)
		if filledWidth < 0 {
			filledWidth = 0
		}
		if filledWidth > barWidth {
			filledWidth = barWidth
		}
		emptyWidth := barWidth - filledWidth

		filledPart := ui.BarFilledStyle.Render(strings.Repeat("█", filledWidth))
		emptyPart := ui.BarEmptyStyle.Render(strings.Repeat("░", emptyWidth))

		bottomStr = fmt.Sprintf("%s%s", filledPart, emptyPart)
	}

	content := lipgloss.JoinVertical(lipgloss.Center, timeStr, "", bottomStr)

	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}
