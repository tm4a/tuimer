package notification

import "os/exec"

func Send(title, message string) {
	_ = exec.Command("notify-send", "-u", "critical", "-a", "Tuimer", title, message).Run()
}
