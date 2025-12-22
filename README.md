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
