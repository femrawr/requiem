package utils

import (
	utils_execute "builder/utils/execute"
	"fmt"
	"os/exec"
)

func RunCommand(dir string, program string, env []string, args ...string) error {
	cmd := exec.Command(program, args...)
	cmd.Dir = dir

	if env != nil {
		cmd.Env = env
	}

	utils_execute.SetCmdHidden(cmd)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v\n%s", err, string(out))
	}

	return nil
}
