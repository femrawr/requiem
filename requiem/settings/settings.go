package settings

import (
	"encoding/json"
	"os"

	"requiem/store"
	"requiem/utils"
)

func LoadSettings() error {
	path, err := utils.CreateADS(store.ExecPath, _CONFIG_ADS_NAME)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			RuntimeSettings = runtimeSettings{}
			return nil
		}

		return err
	}

	return json.Unmarshal(data, &RuntimeSettings)
}

func SaveSettings() error {
	path, err := utils.CreateADS(store.ExecPath, _CONFIG_ADS_NAME)
	if err != nil {
		return err
	}

	data, err := json.Marshal(RuntimeSettings)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0666)
}

func SetSettings(fn func(*runtimeSettings)) error {
	fn(&RuntimeSettings)
	return SaveSettings()
}
