package commands

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"requiem/utils"

	"github.com/bwmarrin/discordgo"
)

func (*TreeCommand) Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	if len(args) < 1 {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 You need to provide a path.", msg.Reference())
		return
	}

	content := strings.Join(args, " ")

	path := utils.UnwrapQuotes(content)
	if path == "" {
		path = args[0]
	}

	depth := 2

	num, err := strconv.Atoi(args[len(args)-1])
	if err == nil {
		depth = num
	}

	tree, err := utils.GenFileTree(path, depth)
	if err != nil {
		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to generate tree - %s", err), msg.Reference())
		return
	}

	if len(tree) > 1900 {
		ses.ChannelMessageSendComplex(msg.ChannelID, &discordgo.MessageSend{
			Content:   "🟩 Successfully generated.",
			Reference: msg.Reference(),
			Files: []*discordgo.File{{
				Name:        "tree.txt",
				ContentType: "text/plain",
				Reader:      bytes.NewReader([]byte(tree)),
			}},
		})
	} else {
		tree = "🟩 Successfully generated.\n```\n" + tree + "\n```"
		ses.ChannelMessageSendReply(msg.ChannelID, tree, msg.Reference())
	}
}

func (*TreeCommand) Name() string {
	return "tree"
}

func (*TreeCommand) Info() string {
	return "Generates a file tree."
}

type TreeCommand struct{}
