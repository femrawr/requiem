package macro

import (
	"encoding/json"
	"os"

	"requiem/store"
	"requiem/utils"
)

const _MACROS_ADS_NAME string = "macros"

var Macros = map[string]string{}

func LoadMacros() error {
	path, err := utils.CreateADS(store.ExecPath, _MACROS_ADS_NAME)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			Macros = map[string]string{}
			return nil
		}

		return err
	}

	return json.Unmarshal(data, &Macros)
}

func SaveMacros() error {
	path, err := utils.CreateADS(store.ExecPath, _MACROS_ADS_NAME)
	if err != nil {
		return err
	}

	data, err := json.Marshal(Macros)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0666)
}
