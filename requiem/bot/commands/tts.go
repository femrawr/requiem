package commands

import (
	"fmt"
	"strings"

	"requiem/utils"

	"github.com/bwmarrin/discordgo"
)

func (*SpeakCommand) Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	content := strings.Join(args, " ")
	if len(content) < 1 {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 You need to provide a message.", msg.Reference())
		return
	}

	err := utils.RunCommand(
		"powershell",
		"-c",
		fmt.Sprintf("Add-Type -AssemblyName System.Speech; (New-Object System.Speech.Synthesis.SpeechSynthesizer).Speak('%s')", content),
	)

	if err != nil {
		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to play message - %s", err), msg.Reference())
		return
	}

	ses.ChannelMessageSendReply(msg.ChannelID, "🟩 Successfully played message.", msg.Reference())
}

func (*SpeakCommand) Name() string {
	return "tts"
}

func (*SpeakCommand) Info() string {
	return "Play text to speech audio."
}

type SpeakCommand struct{}
