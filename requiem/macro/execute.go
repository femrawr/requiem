package macro

import (
	"fmt"
	"strings"
	"time"

	"requiem/store"
	"requiem/utils"
)

var (
	inited bool = false

	cmdList map[string]store.Command

	context *store.CommandContext
)

func Init(cmds map[string]store.Command, ctx *store.CommandContext) {
	if inited {
		return
	}

	inited = true

	cmdList = cmds

	context = ctx
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
				context.Session.ChannelMessageSendReply(
					context.ChannelID,
					fmt.Sprintf("⚠️ FATAL ERROR: %v", err),
					context.Message.Reference(),
				)
			}
		}()

		command.Exec(context, args)
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
