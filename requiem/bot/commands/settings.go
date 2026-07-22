package commands

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"requiem/store"
	"requiem/utils"
)

func (*SettingsCommand) Exec(ctx *store.CommandContext, args []string) {
	if len(args) < 1 {
		ctx.ReplyMsg("🟥 You need to provide a flag.")
		return
	}

	content := strings.Join(args, " ")
	if utils.HasFlag(content, "set") {
		setting := utils.UnwrapQuotes(content)
		if setting == "" {
			ctx.ReplyMsg("🟥 You need to provide a setting wrapped in double quotes.")
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
					ctx.ReplyMsg(fmt.Sprintf("🟥 Invalid bool value: %q", value))
					return
				}

				theValue.Field(i).SetBool(b)
			default:
				ctx.ReplyMsg(fmt.Sprintf("🟥 Unsupported type for setting %q", setting))
				return
			}

			break
		}

		if !found {
			ctx.ReplyMsg(fmt.Sprintf("🟥 Setting %q does not exist.", setting))
			return
		}

		err := store.SaveSettings()
		if err != nil {
			ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to save - %s", err))
			return
		}

		ctx.ReplyMsg(fmt.Sprintf("🟩 Successfully set %q to %q.", setting, value))
		return
	}

	if utils.HasFlag(content, "get") {
		setting := utils.UnwrapQuotes(content)
		if setting == "" {
			ctx.ReplyMsg("🟥 You need to provide a setting wrapped in double quotes.")
			return
		}

		theValue := reflect.ValueOf(store.RuntimeSettings)
		theType := theValue.Type()

		for i := 0; i < theType.NumField(); i++ {
			field := theType.Field(i)
			if field.Tag.Get("json") != setting {
				continue
			}

			ctx.ReplyMsg(fmt.Sprintf("%q = `%v`", setting, theValue.Field(i).Interface()))
			return
		}

		ctx.ReplyMsg(fmt.Sprintf("🟥 Setting %q does not exist.", setting))
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

		ctx.ReplyMsg(fmt.Sprintf("Settings:\n```\n%s\n```", settings.String()))
		return
	}

	ctx.ReplyMsg("🟥 Invalid flag.")
}

func (*SettingsCommand) Name() string {
	return "settings"
}

func (*SettingsCommand) Info() string {
	return "Edit runtime settings."
}

type SettingsCommand struct{}
