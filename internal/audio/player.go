package audio

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/gopxl/beep/v2/vorbis"
	"github.com/gopxl/beep/v2/wav"
)

var (
	initOnce   sync.Once
	initErr    error
	deviceRate beep.SampleRate
)

func PlayAlarm(stopChan chan bool) {
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

	initOnce.Do(func() {
		deviceRate = format.SampleRate
		initErr = speaker.Init(deviceRate, deviceRate.N(time.Second/10))
	})
	if initErr != nil {
		return
	}

	loop, err := beep.Loop2(streamer)
	if err != nil {
		return
	}
	if format.SampleRate != deviceRate {
		loop = beep.Resample(4, format.SampleRate, deviceRate, loop)
	}
	ctrl := &beep.Ctrl{Streamer: loop, Paused: false}
	speaker.Play(ctrl)

	<-stopChan
	speaker.Lock()
	ctrl.Paused = true
	speaker.Unlock()
}
