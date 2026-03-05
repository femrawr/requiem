package commands

import (
	"fmt"

	"requiem/funcs"

	"github.com/bwmarrin/discordgo"
)

func (*ScreenshotCommand) Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	ss, err := funcs.TakeScreenshot()
	if err != nil {
		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to take screenshot - %s", err), msg.Reference())
		return
	}

	ses.ChannelMessageSendComplex(msg.ChannelID, &discordgo.MessageSend{
		Reference: msg.Reference(),
		Files: []*discordgo.File{{
			Name:   "ss.jpg",
			Reader: ss,
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
