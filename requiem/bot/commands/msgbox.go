package commands

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"requiem/store"
)

func (*NotifCommand) Exec(ctx *store.CommandContext, args []string) {
	content := strings.Join(args, " ")
	if len(content) < 1 {
		ctx.ReplyMsg("🟥 You need to provide a message.")
		return
	}

	pointer, err := syscall.UTF16PtrFromString(content)
	if err != nil {
		ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to convert message - %s", err))
		return
	}

	initial, _ := ctx.ReplyMsg("🟩 Successfully sent messagebox.")

	go func() {
		ret, _, err := store.MessageBox.Call(
			uintptr(0),
			uintptr(unsafe.Pointer(pointer)),
			uintptr(0),
			uintptr(0x00000000|0x00040000|0x00001000),
		)

		if ret == 0 {
			ctx.DeleteMsg(initial.ID)
			ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to create messagebox - %s", err))
			return
		}

		ctx.ReplyMsg("Messagebox acknowledgeded.")
	}()
}

func (*NotifCommand) Name() string {
	return "msgbox"
}

func (*NotifCommand) Info() string {
	return "Displays a messagebox."
}

type NotifCommand struct{}
