package store

import "os"

var (
	IsAdmin bool

	ExecPath string

	HomePath string
)

func InitState() {
	admin, _, _ := AmIAdmin.Call()
	IsAdmin = admin != 0

	execPath, err := os.Executable()
	if err == nil {
		ExecPath = execPath
	}

	homePath, err := os.UserHomeDir()
	if err == nil {
		HomePath = homePath
	}
}
