package buttons

import (
	"requiem/store"
	"requiem/utils/discord"

	"github.com/bwmarrin/discordgo"
)

func (*TestButton) Exec(ses *discordgo.Session, itr *discordgo.InteractionCreate, cmd store.Command) {
	messageCreate, _ := discord.GetMessageCreateFromInteraction(ses, itr)

	context := &store.CommandContext{
		Session:     ses,
		Message:     messageCreate,
		ChannelID:   messageCreate.ChannelID,
		Content:     messageCreate.Content,
		Attachments: messageCreate.Attachments,
		Author:      messageCreate.Author,
	}

	cmd.Exec(context, []string{})
}

func (*TestButton) Text() string {
	return "The test button"
}

func (*TestButton) Iden() string {
	return "ping.test"
}

type TestButton struct{}
