package commands

import (
	"fmt"
	"os"
	"strings"

	"requiem/macro"
	"requiem/store"
	"requiem/utils"
	"requiem/utils/discord"
)

func (*MacroCommand) Exec(ctx *store.CommandContext, args []string) {
	if len(args) < 1 {
		ctx.ReplyMsg("🟥 You need to provide a flag or a macro.")
		return
	}

	content := strings.Join(args, " ")
	if utils.HasFlag(content, "register") {
		name := utils.UnwrapQuotes(content)
		if name == "" {
			ctx.ReplyMsg("🟥 You need to provide a name wrapped in double quotes.")
			return
		}

		urls := discord.GetUrls(ctx)
		if len(urls) == 0 {
			ctx.ReplyMsg("🟥 Failed to find any urls.")
			return
		}

		path, err := utils.DownloadFile(urls[0], "")
		if err != nil {
			ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to download - %s", err))
			return
		}

		defer os.Remove(path)

		err = macro.ValidateFile(path)
		if err != nil {
			ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to validate file - %s", err))
			return
		}

		_, exists := macro.Macros[name]
		if exists {
			ctx.ReplyMsg(fmt.Sprintf("🟥 Macro %q already exists.", name))
			return
		}

		parsed, err := macro.ParseMacro(path)
		if err != nil {
			ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to parse macro - %s", err))
			return
		}

		macro.Macros[name] = parsed.Encode()

		err = macro.SaveMacros()
		if err != nil {
			ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to save macros - %s", err))
		}

		ctx.ReplyMsg(fmt.Sprintf("🟩 Successfully registerd macro %q", name))
		return
	}

	if utils.HasFlag(content, "unregister") {
		name := utils.UnwrapQuotes(content)
		if name == "" {
			ctx.ReplyMsg("🟥 You need to provide a name wrapped in double quotes.")
			return
		}

		_, exists := macro.Macros[name]
		if !exists {
			ctx.ReplyMsg("🟥 This macro does not exist.")
			return
		}

		delete(macro.Macros, name)

		err := macro.SaveMacros()
		if err != nil {
			ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to save macros - %s", err))
			return
		}

		ctx.ReplyMsg(fmt.Sprintf("🟩 Successfully unregistered macro %q.", name))
		return
	}

	if utils.HasFlag(content, "list") {
		if len(macro.Macros) == 0 {
			ctx.ReplyMsg("No macros registered.")
			return
		}

		var macros strings.Builder
		for name := range macro.Macros {
			macros.WriteString(name + "\n")
		}

		ctx.ReplyMsg(fmt.Sprintf("Macros:\n```\n%s\n```", macros.String()))
		return
	}

	theMacro, exists := macro.Macros[args[0]]
	if !exists {
		ctx.ReplyMsg("🟥 This macro does not exist.")
		return
	}

	err := macro.RunMacro(theMacro)
	if err != nil {
		ctx.ReplyMsg(fmt.Sprintf("🟥 Failed to run macro - %s", err))
		return
	}

	ctx.ReplyMsg("🟩 Successfully ran macro.")
}

func (*MacroCommand) Name() string {
	return "macro"
}

func (*MacroCommand) Info() string {
	return "Set macros to execute multiple commands at once and more."
}

type MacroCommand struct{}
