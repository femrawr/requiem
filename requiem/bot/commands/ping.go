package commands

import (
	"github.com/bwmarrin/discordgo"
)

func (*PingCommand) Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	ses.ChannelMessageSendReply(
		msg.ChannelID,
		"Pong.",
		msg.Reference(),
	)
}

func (*PingCommand) Name() string {
	return "ping"
}

func (*PingCommand) Info() string {
	return "Test command."
}

type PingCommand struct{}
