package macro

import (
	"fmt"
	"strings"
	"time"

	"requiem/store"
	"requiem/utils"

	"github.com/bwmarrin/discordgo"
)

var (
	inited bool = false

	cmdList map[string]store.Command

	session *discordgo.Session
	message *discordgo.MessageCreate
)

func Init(cmds map[string]store.Command, ses *discordgo.Session, msg *discordgo.MessageCreate) {
	if inited {
		return
	}

	cmdList = cmds
	session = ses
	message = msg

	inited = true
}

func RunMacro(macro string) error {
	chunks := decodeMacro(macro)
	for _, parts := range chunks {
		if len(parts) < 2 {
			return fmt.Errorf("invalid chunk: %s", strings.Join(parts, " "))
		}

		symbol, value := parts[0], parts[1]
		var args []string

		if len(parts) == 3 {
			args = strings.Split(parts[2], ";")
		}

		err := runLine(symbol, value, args)
		if err != nil {
			return err
		}
	}

	return nil
}

func runLine(symbol string, value string, args []string) error {
	switch symbol {
	case "0": // CMD
		return execDiscordCommand(value, args)
	case "1": // WAIT
		secs, found := utils.FindNumber(value)
		if !found {
			return fmt.Errorf("invalid wait value: %q", value)
		}

		time.Sleep(time.Duration(secs) * time.Second)
		return nil
	default:
		return fmt.Errorf("invalid symbol id: %q", symbol)
	}
}

func execDiscordCommand(cmd string, args []string) error {
	command, exists := cmdList[cmd]
	if !exists {
		return fmt.Errorf("command %q does not exist", cmd)
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				session.ChannelMessageSendReply(
					message.ChannelID,
					fmt.Sprintf("⚠️ FATAL ERROR: %v", err),
					message.Reference(),
				)
			}
		}()

		command.Exec(session, message, args)
	}()

	return nil
}

func decodeMacro(macro string) [][]string {
	chunks := strings.Split(macro, "+")

	result := make([][]string, 0, len(chunks))
	for _, chunk := range chunks {
		result = append(result, strings.SplitN(chunk, ".", 3))
	}

	return result
}
