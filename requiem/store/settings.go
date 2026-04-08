package store

import (
	"encoding/json"
	"os"
)

const CONFIG_ADS_NAME string = "c"

var RuntimeSettings Settings

type Settings struct {
	AudioDisableInputsUntilFinished bool `json:"audio_disable_inputs_until_finished"`
	AudioUnmuteBeforePlay           bool `json:"audio_unmute_before_play"`
	AudioMaxVolumeBeforePlay        bool `json:"audio_max_volume_before_play"`

	JumpscareMaxBrightnessBefore        bool `json:"jumpscare_max_brightness_before"`
	JumpscareDisableInputsUntilFinished bool `json:"jumpscare_disable_inputs_until_finished"`
}

func LoadSettings() error {
	data, err := os.ReadFile(ExecPath + ":" + CONFIG_ADS_NAME)
	if err != nil {
		if os.IsNotExist(err) {
			RuntimeSettings = Settings{}
			return nil
		}

		return err
	}

	return json.Unmarshal(data, &RuntimeSettings)
}

func SaveSettings() error {
	data, err := json.Marshal(RuntimeSettings)
	if err != nil {
		return err
	}

	return os.WriteFile(ExecPath+":"+CONFIG_ADS_NAME, data, 0666)
}

func SetSettings(fn func(*Settings)) error {
	fn(&RuntimeSettings)
	return SaveSettings()
}
