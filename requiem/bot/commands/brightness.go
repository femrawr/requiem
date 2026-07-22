package commands

import (
	"fmt"
	"strconv"
	"strings"

	"requiem/store"
	"requiem/utils"
)

func (*LightCommand) Exec(ctx *store.CommandContext, args []string) {
	content := strings.Join(args, " ")
	if len(content) < 1 {
		ctx.ReplyMsg("🟥 You need to provide a brightness level.")
		return
	}

	level, err := strconv.Atoi(args[0])
	if err != nil {
		ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to resolve brightness - %s", err))
		return
	}

	if level < 0 || level > 100 {
		ctx.ReplyMsg("🟥 Light level needs to be between 1 and 100.")
		return
	}

	err = utils.RunCommand(
		"powershell",
		"-c",
		fmt.Sprintf("(Get-WmiObject -Namespace root/WMI -Class WmiMonitorBrightnessMethods).WmiSetBrightness(1, %d)", level),
	)

	if err != nil {
		ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to set brightness - %s", err))
		return
	}

	ctx.ReplyMsg("🟩 Successfully set brightness.")
}

func (*LightCommand) Name() string {
	return "brightness"
}

func (*LightCommand) Info() string {
	return "Sets the device brightness."
}

type LightCommand struct{}
