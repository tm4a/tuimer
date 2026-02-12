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
		if m.Width < 40 {
			barHeight := int(float64(m.Height) * 0.95)
			if barHeight < 5 {
				barHeight = 5
			}

			filledHeight := int(float64(barHeight) * m.Percent)
			if filledHeight < 0 {
				filledHeight = 0
			}
			if filledHeight > barHeight {
				filledHeight = barHeight
			}
			emptyHeight := barHeight - filledHeight

			var lines []string
			for range emptyHeight {
				lines = append(lines, ui.BarEmptyStyle.Render("░"))
			}
			for range filledHeight {
				lines = append(lines, ui.BarFilledStyle.Render("█"))
			}
			bottomStr = strings.Join(lines, "\n")
		} else {
			barWidth := int(float64(m.Width) * 0.95)
			if barWidth < 10 {
				barWidth = 10
			}

			filledWidth := int(float64(barWidth) * m.Percent)
			if filledWidth < 0 {
				filledWidth = 0
			}
			if filledWidth > barWidth {
				filledWidth = barWidth
			}

			if m.Height < 4 {
				textWidth := lipgloss.Width(timeStr)
				textStart := max(0, (barWidth-textWidth-2)/2)
				textEnd := textStart + textWidth + 2

				leftFilled := min(filledWidth, textStart)
				leftEmpty := textStart - leftFilled
				rightFilled := max(0, filledWidth-textEnd)
				rightEmpty := max(0, barWidth-textEnd-rightFilled)

				left := ui.BarFilledStyle.Render(strings.Repeat("█", leftFilled)) +
					ui.BarEmptyStyle.Render(strings.Repeat("░", leftEmpty))
				right := ui.BarFilledStyle.Render(strings.Repeat("█", rightFilled)) +
					ui.BarEmptyStyle.Render(strings.Repeat("░", rightEmpty))

				bottomStr = left + " " + timeStr + " " + right
			} else {
				emptyWidth := barWidth - filledWidth
				filledPart := ui.BarFilledStyle.Render(strings.Repeat("█", filledWidth))
				emptyPart := ui.BarEmptyStyle.Render(strings.Repeat("░", emptyWidth))
				bottomStr = fmt.Sprintf("%s%s", filledPart, emptyPart)
			}
		}
	}

	var content string
	if m.Width >= 40 && m.Height < 4 && !m.InputMode && !m.Finished {
		content = bottomStr
	} else if m.Width < 40 && !m.InputMode && !m.Finished && m.Height < 4 {
		content = timeStr
	} else if m.Width < 40 && !m.InputMode && !m.Finished && (m.Width < 10 || m.Height < 6) {
		content = bottomStr
	} else if bottomStr == "" {
		content = timeStr
	} else {
		content = lipgloss.JoinVertical(lipgloss.Center, timeStr, "", bottomStr)
	}

	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}
