package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tm4a/tuimer/internal/timer"
)

type Model struct {
	Width  int
	Height int

	InputBuffer   string
	InputMode     bool
	Duration      int
	TimeLeft      int
	Percent       float64
	Quitting      bool
	Finished      bool
	Paused        bool
	StopSoundChan chan bool
	SoundPlaying  bool
}

func NewModel(startSeconds int) Model {
	m := Model{
		InputBuffer:   "",
		StopSoundChan: make(chan bool),
	}

	if startSeconds > 0 {
		m.InputMode = false
		m.Duration = startSeconds
		m.TimeLeft = startSeconds
		m.Percent = 1.0
	} else {
		m.InputMode = true
	}

	return m
}

func (m Model) Init() tea.Cmd {
	if !m.InputMode {
		return timer.Tick()
	}
	return nil
}
