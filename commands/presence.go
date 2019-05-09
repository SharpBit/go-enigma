package commands

import (
	"fmt"
	"strings"
)

func presence(ctx *Context) {
	args := ctx.Args

	switch args[0] {
	case "play":
		ctx.Session.UpdateStatus(0, strings.Join(args[1:], " "))
		ctx.Send("Set `play` status to `" + strings.Join(args[1:], " ") + "`")
	case "listen":
		ctx.Session.UpdateListeningStatus(strings.Join(args[1:], " "))
		ctx.Send("Set `listen` status to `" + strings.Join(args[1:], " ") + "`")
	case "stream":
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
	cmd, existing := NewCommand("presence", "Changes the bot's presence.")
	if existing {
		fmt.Println("error: command presence already exists")
		return
	}
	cmd.Dev = true
	cmd.Run = presence
	cmd.Aliases = []string{}
	RegisterCommand(cmd)
}
