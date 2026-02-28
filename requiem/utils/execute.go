package utils

import (
	"os/exec"
	"syscall"
)

func RunCommand(program string, args ...string) error {
	cmd := exec.Command(program, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}

	return cmd.Run()
}

func StartCommand(program string, args ...string) *exec.Cmd {
	cmd := exec.Command(program, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}

	return cmd
}
