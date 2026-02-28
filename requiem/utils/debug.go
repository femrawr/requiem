package utils

import (
	"os"
	"path"

	"requiem/store"
)

const LOG_FILE_NAME string = "requiem debug.txt"

var logFilePath string = ""

func DebugLog(log string) {
	if !store.DEBUG_MODE {
		return
	}

	path := getLogFilePath()
	if path == "" {
		return
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return
	}

	defer file.Close()

	file.WriteString(log + "\n")
}

func getLogFilePath() string {
	if logFilePath != "" {
		return logFilePath
	}

	logFilePath = path.Join(store.HomePath, "Desktop", LOG_FILE_NAME)
	return logFilePath
}
