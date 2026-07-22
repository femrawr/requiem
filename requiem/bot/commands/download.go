package commands

import (
	"fmt"
	"os"
	"strings"

	"requiem/store"
	"requiem/utils"
	"requiem/utils/discord"
)

func (*DownloadCommand) Exec(ctx *store.CommandContext, args []string) {
	urls := discord.GetUrls(ctx)
	if len(urls) == 0 {
		ctx.ReplyMsg("🟥 Failed to find any urls.")
		return
	}

	initial, _ := ctx.ReplyMsg(fmt.Sprintf("Downloading %d files...", len(urls)))

	content := strings.Join(args, " ")
	outPath := utils.UnwrapQuotes(content)

	info, err := os.Stat(outPath)
	if err != nil || !info.IsDir() {
		outPath = ""
	}

	paths, err := utils.DownloadFiles(urls, outPath)
	if err != nil {
		ctx.DeleteMsg(initial.ID)
		ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to download - %s", err))
		return
	}

	var downloads strings.Builder
	downloads.WriteString("🟩 Successfully downloaded.\n```\n")

	for _, path := range paths {
		downloads.WriteString(path + "\n")
	}

	downloads.WriteString("```")

	ctx.DeleteMsg(initial.ID)
	ctx.ReplyMsg(downloads.String())
}

func (*DownloadCommand) Name() string {
	return "download"
}

func (*DownloadCommand) Info() string {
	return "Downloads a file from the server to the device."
}

type DownloadCommand struct{}
