package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"requiem/store"
	"requiem/utils"
)

func (*SiteCommand) Exec(ctx *store.CommandContext, args []string) {
	if len(args) < 2 {
		ctx.ReplyMsg("🟥 You need to provide a flag and a website.")
		return
	}

	content := strings.Join(args, " ")

	list := utils.HasFlag(content, "list")

	site := utils.UnwrapQuotes(content)
	if site == "" && !list {
		ctx.ReplyMsg("🟥 You need to wrap the website in double quotes.")
		return
	}

	path := filepath.Join(os.Getenv("SYSTEMROOT"), "System32", "drivers", "etc", "hosts")

	file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to open file - %s", err))
		return
	}

	defer file.Close()

	site = strings.Replace(site, "https://", "", 1)
	site = strings.Replace(site, "http://", "", 1)
	site = fmt.Sprintf("0.0.0.0 %s", site)

	if utils.HasFlag(content, "block") {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if strings.TrimSpace(scanner.Text()) != site {
				continue
			}

			ctx.ReplyMsg("🟥 This site is already blocked.")
			return
		}

		_, err = file.WriteString("\n" + site)
		if err != nil {
			ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to write file - %s", err))
			return
		}

		ctx.ReplyMsg("🟩 Successfully blocked website.")
		return
	}

	if utils.HasFlag(content, "unblock") {
		data, err := os.ReadFile(path)
		if err != nil {
			ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to read file - %s", err))
			return
		}

		found := false
		newLines := []string{}

		for line := range strings.SplitSeq(string(data), "\n") {
			if strings.TrimSpace(line) == site {
				found = true
				continue
			}

			newLines = append(newLines, line)
		}

		if !found {
			ctx.ReplyMsg("🟥 This site is not blocked.")
			return
		}

		err = os.WriteFile(path, []byte(strings.Join(newLines, "\n")), 0666)
		if err != nil {
			ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to write file - %s", err))
			return
		}

		ctx.ReplyMsg("🟩 Successfully unblocked website.")
		return
	}

	if list {
		data, err := os.ReadFile(path)
		if err != nil {
			ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to read file - %s", err))
			return
		}

		var sites []string
		for line := range strings.SplitSeq(string(data), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			parts := strings.Fields(line)
			if len(parts) < 2 {
				continue
			}

			sites = append(sites, parts[1])
		}

		if len(sites) == 0 {
			ctx.ReplyMsg("There are no blocked sites.")
			return
		}

		ctx.ReplyMsg("Blocked sites:\n```\n" + strings.Join(sites, "\n") + "```")
		return
	}

	ctx.ReplyMsg("🟥 Invalid flag.")
}

func (*SiteCommand) Name() string {
	return "site"
}

func (*SiteCommand) Info() string {
	return "Block or unblock websites."
}

type SiteCommand struct{}
