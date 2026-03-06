package macro

import (
	"fmt"
	"strings"

	"requiem/store"

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

func runLine(symbol, value string, args []string) error {
	switch symbol {
	case "0": // CMD
		return execDiscordCommand(value, args)
	default:
		return fmt.Errorf("invalid symbol id: \"%s\"", symbol)
	}
}

func execDiscordCommand(cmd string, args []string) error {
	command, exists := cmdList[cmd]
	if !exists {
		return fmt.Errorf("command \"%s\" does not exist", cmd)
	}

	go command.Exec(session, message, args)

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
