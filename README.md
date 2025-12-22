# Tuimer

**Tuimer** is a minimal, beautiful, and highly functional terminal timer written in Go.

## ✨ Features

- **Minimalist UI:** Clean interface built with Bubble Tea & Lipgloss.
- **Adaptive:** Resizes gracefully with your terminal window.
- **Smart Controls:** 
  - Type numbers like a microwave (e.g., `130` becomes `01:30`).
  - Space to **Pause/Resume**.
  - **Pomodoro presets** (hidden keys: `p`, `s`, `l`).
- **Notifications:** Desktop notifications via `notify-send` when time is up.
- **Sound:** Plays alarm sound from your config directory.
- **CLI Arguments:** Quick start support (e.g., `tuimer 25m`).

## 📦 Installation

### From Source (Go required)
```bash
go install github.com/tm4a/tuimer@latest
```
### Arch Linux (AUR)
```bash
yay -S tuimer
```

## 🚀 Usage

### Interactive Mode
Just run the app:
```bash
tuimer
```
Type the time (e.g., `10` for 10s, `130` for 1m30s) and hit **Enter**.

### Quick Start
```bash
tuimer 10m    # Start a 10-minute timer
tuimer 1h30m  # Start a 1.5-hour timer
```

### Controls
| Key | Action |
| :--- | :--- |
| **0-9** | Input time (HH:MM:SS) |
| **Enter** | Start Timer / Reset |
| **Space** | Pause / Resume |
| **q** | Reset to input / Quit |
| **Ctrl+C** | Force Quit |

### ⚙️ Configuration (Sound)
To enable the alarm sound, place any `.mp3`, `.wav`, or `.ogg` file in the config directory:
```bash
mkdir -p ~/.config/tuimer
cp /path/to/your/alarm.mp3 ~/.config/tuimer/
```
