package komorebic

import (
	"os/exec"
	"syscall"
)

func Exec(args []string) (string, error) {
	cmd := exec.Command("komorebic.exe", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.Output()
	return string(output), err
}
