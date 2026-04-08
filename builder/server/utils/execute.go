package utils

import (
	"os/exec"
	"syscall"
)

func RunCommand(dir string, program string, args ...string) error {
	cmd := exec.Command(program, args...)

	cmd.Dir = dir
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}

	return cmd.Run()
}
