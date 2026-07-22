package commands

import "requiem/store"

func (*PingCommand) Exec(ctx *store.CommandContext, args []string) {
	ctx.ReplyMsg("Pong.")
}

func (*PingCommand) Name() string {
	return "ping"
}

func (*PingCommand) Info() string {
	return "Test command."
}

type PingCommand struct{}
