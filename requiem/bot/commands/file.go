package commands

import (
	"fmt"
	"os"
	"strings"

	"requiem/utils"

	"github.com/bwmarrin/discordgo"
)

func (*FileCommand) Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	if len(args) < 2 {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 You need to provide a flag and a path.", msg.Reference())
		return
	}

	content := strings.Join(args, " ")

	path := utils.UnwrapQuotes(content)
	if path == "" {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 You need to wrap the path in double quotes.", msg.Reference())
		return
	}

	if utils.HasFlag(content, "delete") {
		err := os.RemoveAll(path)
		if err != nil {
			ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to delete - %s", err), msg.Reference())
			return
		}

		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟩 Successfully deleted - %s", path), msg.Reference())
	} else {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Invalid flag.", msg.Reference())
		return
	}
}

func (*FileCommand) Name() string {
	return "file"
}

func (*FileCommand) Info() string {
	return "Do things with the files on the deivce."
}

type FileCommand struct{}
