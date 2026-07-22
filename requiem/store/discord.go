package store

import "github.com/bwmarrin/discordgo"

type CommandContext struct {
	Session     *discordgo.Session
	Message     *discordgo.MessageCreate
	ChannelID   string
	Content     string
	Attachments []*discordgo.MessageAttachment
	Author      *discordgo.User
}

func (ctx *CommandContext) GetReferenceMsg() *discordgo.Message {
	return ctx.Message.ReferencedMessage
}

func (ctx *CommandContext) ReplyMsg(text string) (*discordgo.Message, error) {
	return ctx.Session.ChannelMessageSendReply(ctx.ChannelID, text, ctx.Message.Reference())
}

func (ctx *CommandContext) SendComplexMsg(data *discordgo.MessageSend) (*discordgo.Message, error) {
	return ctx.Session.ChannelMessageSendComplex(ctx.ChannelID, data)
}

func (ctx *CommandContext) DeleteMsg(id string) error {
	return ctx.Session.ChannelMessageDelete(ctx.ChannelID, id)
}

func (ctx *CommandContext) EditMsg(text string, id string) (*discordgo.Message, error) {
	return ctx.Session.ChannelMessageEdit(ctx.ChannelID, id, text)
}

type Command interface {
	Exec(ctx *CommandContext, args []string)
	Name() string
	Info() string
}

type Button interface {
	Exec(ses *discordgo.Session, itr *discordgo.InteractionCreate, cmd Command)
	Text() string
	Iden() string
}
