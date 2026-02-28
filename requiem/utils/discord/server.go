package discord

import (
	"strings"

	"requiem/funcs"
	"requiem/store"

	"github.com/bwmarrin/discordgo"
)

const DEFAULT_CATEGORY_NAME string = "string2"

func FindCategory(ses *discordgo.Session) string {
	channels, err := ses.GuildChannels(store.SERVER_ID)
	if err != nil {
		return ""
	}

	for _, channel := range channels {
		if channel.Type != discordgo.ChannelTypeGuildCategory {
			continue
		}

		if strings.ToLower(channel.Name) != DEFAULT_CATEGORY_NAME {
			continue
		}

		return channel.ID
	}

	channel, err := ses.GuildChannelCreateComplex(store.SERVER_ID, discordgo.GuildChannelCreateData{
		Name: DEFAULT_CATEGORY_NAME,
		Type: discordgo.ChannelTypeGuildCategory,
	})

	if err != nil {
		return ""
	}

	return channel.ID
}

// the 2nd return is if the channel was newly created
func FindChannel(ses *discordgo.Session, categoryID string) (string, bool) {
	channels, err := ses.GuildChannels(store.SERVER_ID)
	if err != nil {
		return "", false
	}

	fingerprint := funcs.GenFingerprint()

	for _, channel := range channels {
		if channel.Topic != fingerprint || channel.Name != fingerprint {
			continue
		}

		if channel.Type != discordgo.ChannelTypeGuildText {
			continue
		}

		if channel.ParentID != categoryID {
			continue
		}

		return channel.ID, false
	}

	channel, err := ses.GuildChannelCreateComplex(store.SERVER_ID, discordgo.GuildChannelCreateData{
		Name:     fingerprint,
		Topic:    fingerprint,
		Type:     discordgo.ChannelTypeGuildText,
		ParentID: categoryID,
	})

	if err != nil {
		return "", false
	}

	return channel.ID, true
}
