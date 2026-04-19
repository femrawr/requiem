package utils

import (
	"os/exec"
	"strings"
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

func GetCommandOutput(program string, args ...string) (string, error) {
	cmd := exec.Command(program, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}

	out, err := cmd.Output()
	return strings.TrimSpace(string(out)), err
}
