package commands

import (
	"fmt"
	"requiem/utils"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (*LightCommand) Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	content := strings.Join(args, " ")
	if len(content) < 1 {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 You need to provide a brightness level.", msg.Reference())
		return
	}

	level, err := strconv.Atoi(args[0])
	if err != nil {
		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to resolve brightness - %s", err), msg.Reference())
		return
	}

	if level < 0 || level > 100 {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Light level needs to be between 1 and 100.", msg.Reference())
		return
	}

	err = utils.RunCommand(
		"powershell",
		"-c",
		fmt.Sprintf("(Get-WmiObject -Namespace root/WMI -Class WmiMonitorBrightnessMethods).WmiSetBrightness(1, %d)", level),
	)

	if err != nil {
		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to set brightness - %s", err), msg.Reference())
		return
	}

	ses.ChannelMessageSendReply(msg.ChannelID, "🟩 Successfully set brightness.", msg.Reference())
}

func (*LightCommand) Name() string {
	return "brightness"
}

func (*LightCommand) Info() string {
	return "Sets the device brightness."
}

type LightCommand struct{}
