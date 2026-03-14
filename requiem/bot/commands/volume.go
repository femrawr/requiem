package commands

import (
	"fmt"
	"strconv"
	"strings"

	"requiem/funcs"
	"requiem/utils"

	"github.com/bwmarrin/discordgo"
)

func (*VolumeCommand) Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	content := strings.Join(args, " ")
	if len(content) < 1 {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 You need to provide a volume level.", msg.Reference())
		return
	}

	if utils.HasFlag(content, "toggle") {
		err := utils.RunCommand("powershell", "-c", "(New-Object -ComObject WScript.Shell).SendKeys([char]173)")
		if err != nil {
			ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Failed to toggle mute.", msg.Reference())
			return
		}

		ses.ChannelMessageSendReply(msg.ChannelID, "🟩 Successfully toggled mute.", msg.Reference())
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

	err = funcs.SetVolume(float32(level) / 100.0)
	if err == nil {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟩 Successfully set volume.", msg.Reference())
	} else {
		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to set volume - %s", err), msg.Reference())
	}
}

func (*VolumeCommand) Name() string {
	return "volume"
}

func (*VolumeCommand) Info() string {
	return "Sets the device volume."
}

type VolumeCommand struct{}
