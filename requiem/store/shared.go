package store

import "github.com/bwmarrin/discordgo"

type Command interface {
	Exec(ses *discordgo.Session, msg *discordgo.MessageCreate, args []string)
	Name() string
	Info() string
}
