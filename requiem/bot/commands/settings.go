package commands

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"requiem/store"
	"requiem/utils"

	"github.com/bwmarrin/discordgo"
)

func (*SettingsCommand) Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	if len(args) < 1 {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 You need to provide a flag.", msg.Reference())
		return
	}

	content := strings.Join(args, " ")
	if utils.HasFlag(content, "set") {
		setting := utils.UnwrapQuotes(content)
		if setting == "" {
			ses.ChannelMessageSendReply(
				msg.ChannelID,
				"🟥 You need to provide a setting wrapped in double quotes.",
				msg.Reference(),
			)

			return
		}

		value := args[len(args)-1]
		found := false

		theValue := reflect.ValueOf(&store.RuntimeSettings).Elem()
		theType := theValue.Type()

		for i := 0; i < theType.NumField(); i++ {
			field := theType.Field(i)
			if field.Tag.Get("json") != setting {
				continue
			}

			found = true

			switch field.Type.Kind() {
			case reflect.String:
				theValue.Field(i).SetString(value)
			case reflect.Bool:
				b, err := strconv.ParseBool(value)
				if err != nil {
					ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Invalid bool value: %q", value), msg.Reference())
					return
				}

				theValue.Field(i).SetBool(b)
			default:
				ses.ChannelMessageSendReply(
					msg.ChannelID,
					fmt.Sprintf("🟥 Unsupported type for setting %q", setting),
					msg.Reference(),
				)

				return
			}

			break
		}

		if !found {
			ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Setting %q does not exist.", setting), msg.Reference())
			return
		}

		err := store.SaveSettings()
		if err != nil {
			ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to save - %s", err), msg.Reference())
			return
		}

		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟩 Successfully set %q to %q.", setting, value), msg.Reference())
		return
	}

	if utils.HasFlag(content, "get") {
		setting := utils.UnwrapQuotes(content)
		if setting == "" {
			ses.ChannelMessageSendReply(msg.ChannelID, "🟥 You need to provide a setting wrapped in double quotes.", msg.Reference())
			return
		}

		theValue := reflect.ValueOf(store.RuntimeSettings)
		theType := theValue.Type()

		for i := 0; i < theType.NumField(); i++ {
			field := theType.Field(i)
			if field.Tag.Get("json") != setting {
				continue
			}

			ses.ChannelMessageSendReply(
				msg.ChannelID,
				fmt.Sprintf("%q = `%v`", setting, theValue.Field(i).Interface()),
				msg.Reference(),
			)

			return
		}

		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Setting %q does not exist.", setting), msg.Reference())
		return
	}

	if utils.HasFlag(content, "list") {
		theSettings := map[string]any{}

		theValue := reflect.ValueOf(store.RuntimeSettings)
		theType := theValue.Type()

		for i := 0; i < theType.NumField(); i++ {
			field := theType.Field(i)
			theSettings[field.Tag.Get("json")] = theValue.Field(i).Interface()
		}

		var settings strings.Builder
		for key, value := range theSettings {
			fmt.Fprintf(&settings, "%s -> %v\n", key, value)
		}

		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("Settings:\n```\n%s\n```", settings.String()), msg.Reference())
		return
	}

	ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Invalid flag.", msg.Reference())
}

func (*SettingsCommand) Name() string {
	return "settings"
}

func (*SettingsCommand) Info() string {
	return "Edit runtime settings."
}

type SettingsCommand struct{}
