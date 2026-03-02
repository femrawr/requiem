package commands

import (
	"fmt"
	"requiem/funcs"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (*VolumeCommand) Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	content := strings.Join(args, " ")
	if len(content) < 1 {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 You need to provide a volume level.", msg.Reference())
		return
	}

	level, err := strconv.Atoi(args[0])
	if err != nil {
		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to resolve volume - %s", err), msg.Reference())
		return
	}

	if level < 0 || level > 100 {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Volume level needs to be between 1 and 100.", msg.Reference())
		return
	}

	set := funcs.SetVolume(float32(level))
	if set {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟩 Successfully set volume.", msg.Reference())
	} else {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Failed to set volume.", msg.Reference())
	}
}

func (*VolumeCommand) Name() string {
	return "volume"
}

func (*VolumeCommand) Info() string {
	return "Sets the device volume."
}

type VolumeCommand struct{}
