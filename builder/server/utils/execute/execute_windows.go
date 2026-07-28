//go:build windows

package utils_execute

import (
	"os/exec"
	"syscall"
)

func SetCmdHidden(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
}
