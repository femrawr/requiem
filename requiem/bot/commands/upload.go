package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"requiem/utils"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const MAX_FILE_SIZE int64 = 8 * 1024 * 1024

func (*UploadCommand) Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	message := strings.Join(args, " ")
	path := utils.UnwrapQuotes(message)

	info, err := os.Stat(path)
	if err != nil {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Invalid path.", msg.Reference())
		return
	}

	var toUpload string

	initial, _ := ses.ChannelMessageSendReply(msg.ChannelID, "Uploading...", msg.Reference())

	if info.IsDir() {
		zip, err := utils.ZipDir(path)
		if err != nil {
			ses.ChannelMessageDelete(msg.ChannelID, initial.ID)
			ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to zip directory - %s", err), msg.Reference())
			return
		}

		toUpload = zip
		defer os.Remove(zip)
	} else {
		toUpload = path
	}

	info, err = os.Stat(toUpload)
	if err == nil && info.Size() > MAX_FILE_SIZE {
		over := info.Size() - MAX_FILE_SIZE
		ses.ChannelMessageDelete(msg.ChannelID, initial.ID)
		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 File too large by %d bytes.", over), msg.Reference())
		return
	}

	file, err := os.Open(toUpload)
	if err != nil {
		ses.ChannelMessageDelete(msg.ChannelID, initial.ID)
		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to open file - %s", err), msg.Reference())
		return
	}

	defer file.Close()

	ses.ChannelMessageDelete(msg.ChannelID, initial.ID)
	ses.ChannelMessageSendComplex(msg.ChannelID, &discordgo.MessageSend{
		Content:   "🟩 Successfully uploaded.",
		Reference: msg.Reference(),
		Files: []*discordgo.File{
			{
				Name:   filepath.Base(toUpload),
				Reader: file,
			},
		},
	})
}

func (*UploadCommand) Name() string {
	return "upload"
}

func (*UploadCommand) Info() string {
	return "Uploads files from the device to the server."
}

type UploadCommand struct{}
