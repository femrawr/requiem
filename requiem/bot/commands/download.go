package commands

import (
	"fmt"
	"os"
	"strings"

	"requiem/utils"
	"requiem/utils/discord"

	"github.com/bwmarrin/discordgo"
)

func (*DownloadCommand) Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	urls := discord.GetUrls(msg)
	if len(urls) == 0 {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Failed to find any urls.", msg.Reference())
		return
	}

	initial, _ := ses.ChannelMessageSendReply(
		msg.ChannelID,
		fmt.Sprintf("Downloading %d files...", len(urls)),
		msg.Reference(),
	)

	content := strings.Join(args, " ")
	outPath := utils.UnwrapQuotes(content)

	info, err := os.Stat(outPath)
	if err != nil || !info.IsDir() {
		outPath = ""
	}

	paths, err := utils.DownloadFiles(urls, outPath)
	if err != nil {
		ses.ChannelMessageDelete(msg.ChannelID, initial.ID)
		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to download - %s", err), msg.Reference())
		return
	}

	var downloads strings.Builder
	downloads.WriteString("🟩 Successfully downloaded.\n```\n")

	for _, path := range paths {
		downloads.WriteString(path + "\n")
	}

	downloads.WriteString("```")

	ses.ChannelMessageDelete(msg.ChannelID, initial.ID)
	ses.ChannelMessageSendReply(msg.ChannelID, downloads.String(), msg.Reference())
}

func (*DownloadCommand) Name() string {
	return "download"
}

func (*DownloadCommand) Info() string {
	return "Downloads a file from the server to the device."
}

type DownloadCommand struct{}
