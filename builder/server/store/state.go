package store

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	PORT string = "7305"

	DEBUG bool = true
)

var (
	// the root of the whole Requiem project
	Root string

	// the actual dir where Requiem itself is in
	Main string

	// the build tag, it is set in /routes/update_config.go
	Tag string = "none"
)

func InitState() {
	path, err := os.Executable()
	if err != nil {
		fmt.Printf("failed to get the path of the executable - %s", err)
		return
	}

	Root = filepath.Join(path, "..", "..", "..")

	Main = filepath.Join(Root, "requiem")
}
