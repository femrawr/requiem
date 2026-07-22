package discord

import (
	"github.com/bwmarrin/discordgo"
)

func GetMessageCreateFromInteraction(ses *discordgo.Session, itr *discordgo.InteractionCreate) (*discordgo.MessageCreate, error) {
	message, err := ses.ChannelMessage(itr.ChannelID, itr.Message.ID)
	if err != nil {
		return nil, err
	}

	return &discordgo.MessageCreate{
		Message: message,
	}, nil
}
