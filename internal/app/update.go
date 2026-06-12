package app

import (
	"math"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tm4a/tuimer/internal/audio"
	"github.com/tm4a/tuimer/internal/notification"
	"github.com/tm4a/tuimer/internal/timer"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c":
			m.Quitting = true
			if m.SoundPlaying {
				select {
				case m.StopSoundChan <- true:
				default:
				}
			}
			return m, tea.Quit

		case "q":
			if m.InputMode {
				m.Quitting = true
				return m, tea.Quit
			} else {
				if m.SoundPlaying {
					select {
					case m.StopSoundChan <- true:
					default:
					}
					m.SoundPlaying = false
				}
				m.Finished = false
				m.Paused = false
				m.InputMode = true
				m.InputBuffer = ""
				m.Duration = 0
				m.TimeLeft = 0
				m.Percent = 0
				m.TickID++
				return m, nil
			}

		case " ":
			if !m.InputMode && !m.Finished {
				if m.Paused {
					m.Paused = false
					m.Deadline = time.Now().Add(m.PausedRemaining)
					m.TickID++
					return m, timer.Tick(m.TickID, m.Deadline)
				}
				m.Paused = true
				m.PausedRemaining = time.Until(m.Deadline)
				m.TickID++
			}

		case "enter":
			if m.InputMode {
				totalSeconds := timer.ParseAndValidateTime(m.InputBuffer)
				if totalSeconds > 0 {
					return m.startTimer(totalSeconds)
				}
			} else if m.Finished {
				if m.SoundPlaying {
					select {
					case m.StopSoundChan <- true:
					default:
					}
					m.SoundPlaying = false
				}
				m.Finished = false
				m.InputMode = true
				m.InputBuffer = ""
				m.Duration = 0
				m.TimeLeft = 0
				m.Percent = 0
				m.TickID++
				return m, nil
			}

		case "p":
			if m.InputMode {
				return m.startTimer(25 * 60)
			}
		case "s":
			if m.InputMode {
				return m.startTimer(5 * 60)
			}
		case "l":
			if m.InputMode {
				return m.startTimer(15 * 60)
			}

		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			if m.InputMode && len(m.InputBuffer) < 6 {
				m.InputBuffer += msg.String()
			}
		case "backspace":
			if m.InputMode && len(m.InputBuffer) > 0 {
				m.InputBuffer = m.InputBuffer[:len(m.InputBuffer)-1]
			}
		}

	case timer.TickMsg:
		if msg.ID != m.TickID || m.InputMode || m.Paused || m.Finished {
			return m, nil
		}
		remaining := time.Until(m.Deadline)
		if remaining <= 0 {
			m.TimeLeft = 0
			m.Percent = 0
			m.Finished = true
			m.SoundPlaying = true
			go notification.Send("Tuimer", "Time is up!")
			go audio.PlayAlarm(m.StopSoundChan)
			return m, nil
		}
		m.TimeLeft = int(math.Ceil(remaining.Seconds()))
		m.Percent = float64(m.TimeLeft) / float64(m.Duration)
		return m, timer.Tick(m.TickID, m.Deadline)
	}

	return m, nil
}

func (m Model) startTimer(seconds int) (Model, tea.Cmd) {
	m.Duration = seconds
	m.TimeLeft = seconds
	m.InputMode = false
	m.Paused = false
	m.Percent = 1.0
	m.Deadline = time.Now().Add(time.Duration(seconds) * time.Second)
	m.TickID++
	return m, timer.Tick(m.TickID, m.Deadline)
}
