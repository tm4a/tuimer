package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/gopxl/beep/v2/vorbis"
	"github.com/gopxl/beep/v2/wav"
)

// --- STYLES ---
var (
	primaryColor = lipgloss.Color("#36968e")
	subtleColor  = lipgloss.Color("241")
	alertColor   = lipgloss.Color("196")

	timeStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	pausedStyle = lipgloss.NewStyle().
			Foreground(subtleColor).
			Bold(true)

	alarmStyle = lipgloss.NewStyle().
			Foreground(alertColor).
			Bold(true).
			Blink(true)

	barFilledStyle = lipgloss.NewStyle().Foreground(primaryColor)
	barEmptyStyle  = lipgloss.NewStyle().Foreground(subtleColor)
)

// --- MODEL ---
type model struct {
	width  int
	height int

	inputBuffer   string
	inputMode     bool
	duration      int
	timeLeft      int
	percent       float64
	quitting      bool
	finished      bool
	paused        bool
	stopSoundChan chan bool
	soundPlaying  bool
}

type TickMsg time.Time

func initialModel(startSeconds int) model {
	m := model{
		inputBuffer:   "",
		stopSoundChan: make(chan bool),
	}

	if startSeconds > 0 {
		m.inputMode = false
		m.duration = startSeconds
		m.timeLeft = startSeconds
		m.percent = 1.0
	} else {
		m.inputMode = true
	}

	return m
}

func (m model) Init() tea.Cmd {
	if !m.inputMode {
		return tick()
	}
	return nil
}

// --- UPDATE ---
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {

		// 1. HARD QUIT (Ctrl+C always quits)
		case "ctrl+c":
			m.quitting = true
			if m.soundPlaying {
				select {
				case m.stopSoundChan <- true:
				default:
				}
			}
			return m, tea.Quit

		// 2. SOFT QUIT / RESET (q logic)
		case "q":
			if m.inputMode {
				// If we are at the input screen -> Quit program
				m.quitting = true
				return m, tea.Quit
			} else {
				// If timer is running/paused/finished -> Reset to input
				if m.soundPlaying {
					select {
					case m.stopSoundChan <- true:
					default:
					}
					m.soundPlaying = false
				}

				// Reset state
				m.finished = false
				m.paused = false
				m.inputMode = true
				m.inputBuffer = ""
				m.duration = 0
				m.timeLeft = 0
				m.percent = 0
				return m, nil
			}

		case " ":
			if !m.inputMode && !m.finished {
				m.paused = !m.paused
				if !m.paused {
					return m, tick()
				}
			}

		case "enter":
			if m.inputMode {
				totalSeconds, _ := parseAndValidateTime(m.inputBuffer)
				if totalSeconds > 0 {
					return m.startTimer(totalSeconds)
				}
			} else if m.finished {
				// Stop sound and reset (Same logic as 'q' reset)
				if m.soundPlaying {
					select {
					case m.stopSoundChan <- true:
					default:
					}
					m.soundPlaying = false
				}
				m.finished = false
				m.inputMode = true
				m.inputBuffer = ""
				m.duration = 0
				m.timeLeft = 0
				m.percent = 0
				return m, nil
			}

		// Hidden Pomodoro shortcuts
		case "p":
			if m.inputMode {
				return m.startTimer(25 * 60)
			}
		case "s":
			if m.inputMode {
				return m.startTimer(5 * 60)
			}
		case "l":
			if m.inputMode {
				return m.startTimer(15 * 60)
			}

		// Input numbers
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			if m.inputMode && len(m.inputBuffer) < 6 {
				m.inputBuffer += msg.String()
			}
		case "backspace":
			if m.inputMode && len(m.inputBuffer) > 0 {
				m.inputBuffer = m.inputBuffer[:len(m.inputBuffer)-1]
			}
		}

	case TickMsg:
		if m.paused {
			return m, nil
		}
		if !m.inputMode && m.timeLeft > 0 {
			m.timeLeft--
			m.percent = float64(m.timeLeft) / float64(m.duration)
			return m, tick()
		} else if !m.inputMode && m.timeLeft == 0 && !m.finished {
			m.finished = true
			m.soundPlaying = true
			go sendNotification("Tuimer", "Time is up!")
			go playAlarm(m.stopSoundChan)
			return m, nil
		}
	}

	return m, nil
}

func (m model) startTimer(seconds int) (model, tea.Cmd) {
	m.duration = seconds
	m.timeLeft = seconds
	m.inputMode = false
	m.percent = 1.0
	return m, tick()
}

// --- VIEW ---
func (m model) View() string {
	if m.quitting {
		return ""
	}

	// 1. Time Display
	var timeStr string
	if m.inputMode {
		timeStr = formatBuffer(m.inputBuffer)
		timeStr = timeStyle.Render(timeStr)
	} else {
		h := m.timeLeft / 3600
		min := (m.timeLeft % 3600) / 60
		sec := m.timeLeft % 60

		rawTime := fmt.Sprintf("%02d:%02d:%02d", h, min, sec)

		if m.finished {
			timeStr = alarmStyle.Render(rawTime)
		} else if m.paused {
			timeStr = pausedStyle.Render(rawTime + " [PAUSED]")
		} else {
			timeStr = timeStyle.Render(rawTime)
		}
	}

	// 2. Bar
	var bottomStr string

	if !m.inputMode && !m.finished {
		barWidth := int(float64(m.width) * 0.6)
		if barWidth < 10 {
			barWidth = 10
		}
		if barWidth > 100 {
			barWidth = 100
		}

		filledWidth := int(float64(barWidth) * m.percent)
		if filledWidth < 0 {
			filledWidth = 0
		}
		if filledWidth > barWidth {
			filledWidth = barWidth
		}
		emptyWidth := barWidth - filledWidth

		filledPart := barFilledStyle.Render(strings.Repeat("█", filledWidth))
		emptyPart := barEmptyStyle.Render(strings.Repeat("░", emptyWidth))

		bottomStr = fmt.Sprintf("%s%s", filledPart, emptyPart)
	}

	// 3. Assembly
	content := lipgloss.JoinVertical(lipgloss.Center, timeStr, "", bottomStr)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

// --- UTILS ---
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

func sendNotification(title, message string) {
	_ = exec.Command("notify-send", "-u", "critical", "-a", "Tuimer", title, message).Run()
}

func playAlarm(stopChan chan bool) {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	configDir := filepath.Join(home, ".config", "tuimer")
	files, err := os.ReadDir(configDir)
	if err != nil || len(files) == 0 {
		return
	}
	var filePath string
	for _, f := range files {
		ext := strings.ToLower(filepath.Ext(f.Name()))
		if ext == ".mp3" || ext == ".wav" || ext == ".ogg" {
			filePath = filepath.Join(configDir, f.Name())
			break
		}
	}
	if filePath == "" {
		return
	}
	f, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer f.Close()
	var streamer beep.StreamSeekCloser
	var format beep.Format
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".mp3":
		streamer, format, err = mp3.Decode(f)
	case ".wav":
		streamer, format, err = wav.Decode(f)
	case ".ogg":
		streamer, format, err = vorbis.Decode(f)
	}
	if err != nil {
		return
	}
	defer streamer.Close()
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	loop := beep.Loop(-1, streamer)
	ctrl := &beep.Ctrl{Streamer: loop, Paused: false}
	speaker.Play(ctrl)
	<-stopChan
	speaker.Lock()
	ctrl.Paused = true
	speaker.Unlock()
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func formatBuffer(buf string) string {
	padded := buf
	for len(padded) < 6 {
		padded = "0" + padded
	}
	return padded[0:2] + ":" + padded[2:4] + ":" + padded[4:6]
}

func parseAndValidateTime(buf string) (int, string) {
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
	return (h * 3600) + (m * 60) + s, ""
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
			fmt.Println("❌ Error: Invalid time format.")
			fmt.Println("Run 'tuimer --help' for details.")
			os.Exit(1)
		}
		startSeconds = int(d.Seconds())
	}

	p := tea.NewProgram(initialModel(startSeconds), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}
}
