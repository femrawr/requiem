package commands

import (
	"requiem/persistence"

	"github.com/bwmarrin/discordgo"
)

func (*PersistCommand) Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	err := persistence.RunRegistryPersist("", true)
	if err != nil {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Failed to persist (run registry).", msg.Reference())
	}

	err = persistence.SchedularPersist("", true)
	if err != nil {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Failed to persist (schedular).", msg.Reference())
		return
	}

	ses.ChannelMessageSendReply(msg.ChannelID, "🟩 Successfully persisted.", msg.Reference())
}

func (*PersistCommand) Name() string {
	return "persist"
}

func (*PersistCommand) Info() string {
	return "Persists or re-persists Requiem if it was disabled."
}

type PersistCommand struct{}
