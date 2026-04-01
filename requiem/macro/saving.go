package macro

import (
	"encoding/json"
	"os"

	"requiem/store"
)

const MACROS_ADS_NAME string = "m"

var Macros = map[string]string{}

func LoadMacros() error {
	data, err := os.ReadFile(store.ExecPath + ":" + MACROS_ADS_NAME)
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
	data, err := json.Marshal(Macros)
	if err != nil {
		return err
	}

	return os.WriteFile(store.ExecPath+":"+MACROS_ADS_NAME, data, 0666)
}
