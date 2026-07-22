package commands

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"requiem/store"
	"requiem/utils"

	"github.com/bwmarrin/discordgo"
)

func (*TreeCommand) Exec(ctx *store.CommandContext, args []string) {
	if len(args) < 1 {
		ctx.ReplyMsg("🟥 You need to provide a path.")
		return
	}

	content := strings.Join(args, " ")

	path := utils.UnwrapQuotes(content)
	if path == "" {
		path = args[0]
	}

	depth := 2

	num, err := strconv.Atoi(args[len(args)-1])
	if err == nil && num != 0 {
		depth = num
	}

	tree, err := utils.GenFileTree(path, depth)
	if err != nil {
		ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to generate tree - %s", err))
		return
	}

	if len(tree) > 1900 {
		ctx.SendComplexMsg(&discordgo.MessageSend{
			Content:   "🟩 Successfully generated.",
			Reference: ctx.Message.Reference(),
			Files: []*discordgo.File{{
				Name:        "tree.txt",
				ContentType: "text/plain",
				Reader:      bytes.NewReader([]byte(tree)),
			}},
		})
	} else {
		tree = "🟩 Successfully generated.\n```\n" + tree + "\n```"
		ctx.ReplyMsg(tree)
	}
}

func (*TreeCommand) Name() string {
	return "tree"
}

func (*TreeCommand) Info() string {
	return "Generates a file tree."
}

type TreeCommand struct{}
