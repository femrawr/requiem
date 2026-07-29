package store

import (
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

	// the key used to decrypt e2e config, it is set in /routes/get_key.go
	SharedSecret []byte

	// what operations should be dont on the build, it is set in /routes/update_config.go
	Obfuscate    bool = false
	Pack         bool = false
	BuildAs32Bit bool = false
)

func InitState() error {
	path, err := os.Executable()
	if err != nil {
		return err
	}

	Root = filepath.Join(path, "..", "..", "..")
	Main = filepath.Join(Root, "requiem")

	return nil
}
