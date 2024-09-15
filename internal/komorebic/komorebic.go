package komorebic

import (
	"os"
	"os/exec"
)

func Exec(args []string) error {
	cmd := exec.Command("komorebic.exe", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
