package settings

const _CONFIG_ADS_NAME string = "runtime-settings"

var RuntimeSettings runtimeSettings

type runtimeSettings struct {
	AudioDisableInputsUntilFinished bool `json:"audio_disable_inputs_until_finished"`
	AudioUnmuteBeforePlay           bool `json:"audio_unmute_before_play"`
	AudioMaxVolumeBeforePlay        bool `json:"audio_max_volume_before_play"`

	JumpscareMaxBrightnessBefore        bool `json:"jumpscare_max_brightness_before"`
	JumpscareDisableInputsUntilFinished bool `json:"jumpscare_disable_inputs_until_finished"`
}
