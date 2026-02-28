package funcs

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"strings"

	"requiem/utils"
)

func GenFingerprint() string {
	script := `
        (Get-CimInstance Win32_BIOS).SerialNumber +
        (Get-CimInstance Win32_BaseBoard).SerialNumber +
        (Get-CimInstance Win32_ComputerSystemProduct).UUID +
        (Get-CimInstance Win32_DiskDrive | Select -First 1).SerialNumber
    `

	cmd := utils.StartCommand(
		"powershell",
		"-nop",
		"-w hidden",
		"-ep bypass",
		"-c",
		script,
	)

	out, err := cmd.StdoutPipe()
	if err != nil {
		return "1"
	}

	err = cmd.Start()
	if err != nil {
		return "2"
	}

	bytes, err := io.ReadAll(out)
	if err != nil {
		return "3"
	}

	err = cmd.Wait()
	if err != nil {
		return "4"
	}

	trimmed := strings.TrimSpace(string(bytes))
	hash := sha256.Sum256([]byte(trimmed))

	return hex.EncodeToString(hash[:])
}
