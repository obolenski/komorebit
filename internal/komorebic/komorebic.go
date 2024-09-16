package komorebic

import (
	"os/exec"
)

func Exec(args []string) (string, error) {
	cmd := exec.Command("komorebic.exe", args...)
	output, err := cmd.Output()
	return string(output), err
}
