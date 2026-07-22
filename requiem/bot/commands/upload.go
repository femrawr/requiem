package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"requiem/store"
	"requiem/utils"

	"github.com/bwmarrin/discordgo"
)

const _MAX_FILE_SIZE int64 = 8 * 1024 * 1024

func (*UploadCommand) Exec(ctx *store.CommandContext, args []string) {
	message := strings.Join(args, " ")
	path := utils.UnwrapQuotes(message)

	info, err := os.Stat(path)
	if err != nil {
		ctx.ReplyMsg("🟥 Invalid path.")
		return
	}

	var toUpload string

	initial, _ := ctx.ReplyMsg("Uploading...")

	if info.IsDir() {
		zip, err := utils.ZipDir(path)
		if err != nil {
			ctx.DeleteMsg(initial.ID)
			ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to zip directory - %s", err))
			return
		}

		toUpload = zip
		defer os.Remove(zip)
	} else {
		toUpload = path
	}

	info, err = os.Stat(toUpload)
	if err == nil && info.Size() > _MAX_FILE_SIZE {
		over := info.Size() - _MAX_FILE_SIZE
		ctx.DeleteMsg(initial.ID)
		ctx.ReplyMsg(fmt.Sprintf("🟥 File too large by %d bytes.", over))
		return
	}

	file, err := os.Open(toUpload)
	if err != nil {
		ctx.DeleteMsg(initial.ID)
		ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to open file - %s", err))
		return
	}

	defer file.Close()

	ctx.DeleteMsg(initial.ID)
	ctx.SendComplexMsg(&discordgo.MessageSend{
		Content:   "🟩 Successfully uploaded.",
		Reference: ctx.Message.Reference(),
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
