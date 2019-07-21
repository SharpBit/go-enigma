package cogs

import (
	"strings"
)

func presence(ctx *Context) {
	args := ctx.Args
	usage := "Usage: `//play <activity> [url:stream] [message]`"
	if len(args) == 0 {
		ctx.Send(usage)
		return
	}

	switch args[0] {
	case "play":
		if len(args) < 2 {
			ctx.Send(usage)
			return
		}
		ctx.Session.UpdateStatus(0, strings.Join(args[1:], " "))
		ctx.Send("Set `play` status to `" + strings.Join(args[1:], " ") + "`")
	case "listen":
		if len(args) < 2 {
			ctx.Send(usage)
			return
		}
		ctx.Session.UpdateListeningStatus(strings.Join(args[1:], " "))
		ctx.Send("Set `listen` status to `" + strings.Join(args[1:], " ") + "`")
	case "stream":
		if len(args) < 3 {
			ctx.Send(usage)
			return
		}
		CleanURL := strings.Replace(strings.Replace(args[1], "<", "", 1), ">", "", 1)
		ctx.Session.UpdateStreamingStatus(0, strings.Join(args[2:], " "), CleanURL)
		ctx.Send("Set `stream` status to `" + strings.Join(args[2:], " ") + "` at URL <" + CleanURL + ">")
	case "reset":
		ctx.Session.UpdateStatus(0, "")
		ctx.Send("Presence reset.")
	default:
		ctx.Send(args[0] + " is not an option. Select from `play`, `listen`, `stream`, `reset`")
	}
}

func init() {
	cog := NewCog("Owner", "Developer restricted commands", true)
	cog.AddCommand("presence", "Changes the bot's presence", []string{}, presence)
	cog.Load()
}
