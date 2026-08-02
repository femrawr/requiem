package utils

import (
	"fmt"

	"builder/store"

	"github.com/fatih/color"
)

var (
	logInfo  = color.New(color.FgWhite).SprintfFunc()
	logWarn  = color.New(color.FgYellow).SprintfFunc()
	logError = color.New(color.FgRed).SprintfFunc()
)

func LogDebug(msg string, args ...any) {
	if !store.DEBUG {
		return
	}

	if len(args) == 0 {
		fmt.Printf("[DEBUG] %s\n", msg)
	} else {
		fmt.Printf("[DEBUG] %s %s\n", msg, fmt.Sprint(args...))
	}
}

func LogChunk() {
	fmt.Print("\n")
}

func LogInfo(msg string, args ...any) {
	if len(args) == 0 {
		fmt.Print(logInfo("[INFO] %s\n", msg))
	} else {
		fmt.Print(logInfo("[INFO] %s %s\n", msg, fmt.Sprint(args...)))
	}
}

func LogWarning(msg string, args ...any) {
	if len(args) == 0 {
		fmt.Print(logWarn("[WARN] %s\n", msg))
	} else {
		fmt.Print(logWarn("[WARN] %s %s\n", msg, fmt.Sprint(args...)))
	}
}

func LogError(msg string, args ...any) {
	if len(args) == 0 {
		fmt.Print(logError("[ERROR] %s\n", msg))
	} else {
		fmt.Print(logError("[ERROR] %s %s\n", msg, fmt.Sprint(args...)))
	}
}
