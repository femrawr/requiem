package commands

import (
	"fmt"
	"strconv"
	"strings"

	"requiem/funcs"
	"requiem/store"
	"requiem/utils"
)

func (*VolumeCommand) Exec(ctx *store.CommandContext, args []string) {
	content := strings.Join(args, " ")
	if len(content) < 1 {
		ctx.ReplyMsg("🟥 You need to provide a volume level.")
		return
	}

	if utils.HasFlag(content, "mute") {
		err := funcs.SetMuted(true)
		if err != nil {
			ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to mute - %s", err))
			return
		}

		ctx.ReplyMsg("🟩 Successfully muted.")
		return
	}

	if utils.HasFlag(content, "unmute") {
		err := funcs.SetMuted(false)
		if err != nil {
			ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to unmute - %s", err))
			return
		}

		ctx.ReplyMsg("🟩 Successfully unmuted.")
		return
	}

	level, err := strconv.Atoi(args[0])
	if err != nil {
		ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to resolve volume - %s", err))
		return
	}

	if level < 0 || level > 100 {
		ctx.ReplyMsg("🟥 Volume level needs to be between 1 and 100.")
		return
	}

	err = funcs.SetVolume(float32(level) / 100.0)
	if err == nil {
		ctx.ReplyMsg("🟩 Successfully set volume.")
	} else {
		ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to set volume - %s", err))
	}
}

func (*VolumeCommand) Name() string {
	return "volume"
}

func (*VolumeCommand) Info() string {
	return "Sets the device volume."
}

type VolumeCommand struct{}
