package commands

import (
	"fmt"
	"os"
	"regexp"
	"requiem/utils"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (*DownloadCommand) Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	var urls []string

	content := strings.Join(args, " ")

	regex := regexp.MustCompile(`https?://[^\s]+`)

	matches := regex.FindAllString(content, -1)
	urls = append(urls, matches...)

	for _, attachment := range msg.Attachments {
		urls = append(urls, attachment.URL)
	}

	if msg.ReferencedMessage != nil {
		for _, attachment := range msg.ReferencedMessage.Attachments {
			urls = append(urls, attachment.URL)
		}

		matches := regex.FindAllString(msg.ReferencedMessage.Content, -1)
		urls = append(urls, matches...)
	}

	if len(urls) == 0 {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Failed to find any urls.", msg.Reference())
		return
	}

	initial, _ := ses.ChannelMessageSendReply(
		msg.ChannelID,
		fmt.Sprintf("Downloading %d files...", len(urls)),
		msg.Reference(),
	)

	outPath := utils.UnwrapQuotes(content)

	info, err := os.Stat(outPath)
	if err != nil || !info.IsDir() {
		outPath = ""
	}

	paths, err := utils.DownloadFiles(urls, outPath)
	if err != nil {
		ses.ChannelMessageDelete(msg.ChannelID, initial.ID)
		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed download - %s", err), msg.Reference())
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
