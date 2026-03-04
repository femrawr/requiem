package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"requiem/utils"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (*SiteCommand) Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	if len(args) < 2 {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 You need to provide a flag and a website.", msg.Reference())
		return
	}

	content := strings.Join(args, " ")

	site := utils.UnwrapQuotes(content)
	if site == "" {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 You need to wrap the website in double quotes.", msg.Reference())
		return
	}

	path := filepath.Join(os.Getenv("SYSTEMROOT"), "System32", "drivers", "etc", "hosts")

	file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to open file - %s", err), msg.Reference())
		return
	}

	defer file.Close()

	site = strings.Replace(site, "https://", "", 1)
	site = strings.Replace(site, "http://", "", 1)
	site = fmt.Sprintf("127.0.0.1 %s", site)

	if utils.HasFlag(content, "block") {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if strings.TrimSpace(scanner.Text()) != site {
				continue
			}

			ses.ChannelMessageSendReply(msg.ChannelID, "🟥 This site is already blocked.", msg.Reference())
			return
		}

		_, err = file.WriteString("\n" + site)
		if err != nil {
			ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to write file - %s", err), msg.Reference())
			return
		}

		ses.ChannelMessageSendReply(msg.ChannelID, "🟩 Successfully blocked website.", msg.Reference())
	} else if utils.HasFlag(content, "unblock") {
		data, err := os.ReadFile(path)
		if err != nil {
			ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to read file - %s", err), msg.Reference())
			return
		}

		found := false
		newLines := []string{}

		for line := range strings.SplitSeq(string(data), "\n") {
			if strings.TrimSpace(line) == site {
				found = true
				continue
			}

			newLines = append(newLines, line)
		}

		if !found {
			ses.ChannelMessageSendReply(msg.ChannelID, "🟥 This site is not blocked.", msg.Reference())
			return
		}

		err = os.WriteFile(path, []byte(strings.Join(newLines, "\n")), 0666)
		if err != nil {
			ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to write file - %s", err), msg.Reference())
			return
		}

		ses.ChannelMessageSendReply(msg.ChannelID, "🟩 Successfully unblocked website.", msg.Reference())
	} else {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Invalid flag.", msg.Reference())
		return
	}
}

func (*SiteCommand) Name() string {
	return "site"
}

func (*SiteCommand) Info() string {
	return "Block or unblock websites."
}

type SiteCommand struct{}
