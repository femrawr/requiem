package commands

import (
	"fmt"
	"os"
	"strings"

	"requiem/macro"
	"requiem/utils"
	"requiem/utils/discord"

	"github.com/bwmarrin/discordgo"
)

func (*MacroCommand) Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	if len(args) < 1 {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 You need to provide a flag or a macro.", msg.Reference())
		return
	}

	content := strings.Join(args, " ")
	if utils.HasFlag(content, "register") {
		name := utils.UnwrapQuotes(content)
		if name == "" {
			ses.ChannelMessageSendReply(msg.ChannelID, "🟥 You need to provide a name wrapped in double quotes.", msg.Reference())
			return
		}

		urls := discord.GetUrls(msg)
		if len(urls) == 0 {
			ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Failed to find any urls.", msg.Reference())
			return
		}

		path, err := utils.DownloadFile(urls[0], "")
		if err != nil {
			ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to download - %s", err), msg.Reference())
			return
		}

		defer os.Remove(path)

		err = macro.ValidateFile(path)
		if err != nil {
			ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to validate file - %s", err), msg.Reference())
			return
		}

		_, exists := macro.Macros[name]
		if exists {
			ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Macro %q already exists.", name), msg.Reference())
			return
		}

		parsed, err := macro.ParseMacro(path)
		if err != nil {
			ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to parse macro - %s", err), msg.Reference())
			return
		}

		macro.Macros[name] = parsed.Encode()

		err = macro.SaveMacros()
		if err != nil {
			ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to save macros - %s", err), msg.Reference())
		}

		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟩 Successfully registerd macro %q", name), msg.Reference())
		return
	}

	if utils.HasFlag(content, "unregister") {
		name := utils.UnwrapQuotes(content)
		if name == "" {
			ses.ChannelMessageSendReply(msg.ChannelID, "🟥 You need to provide a name wrapped in double quotes.", msg.Reference())
			return
		}

		_, exists := macro.Macros[name]
		if !exists {
			ses.ChannelMessageSendReply(msg.ChannelID, "🟥 This macro does not exist.", msg.Reference())
			return
		}

		delete(macro.Macros, name)

		err := macro.SaveMacros()
		if err != nil {
			ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to save macros - %s", err), msg.Reference())
			return
		}

		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟩 Successfully unregistered macro %q.", name), msg.Reference())
		return
	}

	if utils.HasFlag(content, "list") {
		if len(macro.Macros) == 0 {
			ses.ChannelMessageSendReply(msg.ChannelID, "No macros registered.", msg.Reference())
			return
		}

		var macros strings.Builder
		for name := range macro.Macros {
			macros.WriteString(name + "\n")
		}

		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("Macros:\n```\n%s\n```", macros.String()), msg.Reference())
		return
	}

	theMacro, exists := macro.Macros[args[0]]
	if !exists {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 This macro does not exist.", msg.Reference())
		return
	}

	err := macro.RunMacro(theMacro)
	if err != nil {
		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to run macro - %s", err), msg.Reference())
		return
	}

	ses.ChannelMessageSendReply(msg.ChannelID, "🟩 Successfully ran macro.", msg.Reference())
}

func (*MacroCommand) Name() string {
	return "macro"
}

func (*MacroCommand) Info() string {
	return "Set macros to execute multiple commands at once and more."
}

type MacroCommand struct{}
