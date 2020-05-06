package cogs

import (
	"fmt"
	"strings"

	"github.com/SharpBit/go-enigma/commands"
)

func presence(ctx *commands.Context, ActivityType string, UrlOrMessage ...string) (err error) {
	switch ActivityType {
	case "play":
		if len(UrlOrMessage) < 1 {
			return fmt.Errorf("ArgumentError: You must provide a message to show.")
		}
		message := strings.Join(UrlOrMessage, " ")
		ctx.Session.UpdateStatus(0, message)
		_, err = ctx.Send("Set `play` status to `" + message + "`")
	case "listen":
		if len(UrlOrMessage) < 1 {
			return fmt.Errorf("ArgumentError: You must provide a message to show.")
		}
		message := strings.Join(UrlOrMessage, " ")
		ctx.Session.UpdateListeningStatus(message)
		_, err = ctx.Send("Set `listen` status to `" + message + "`")
	case "stream":
		if len(UrlOrMessage) < 2 {
			return fmt.Errorf("ArgumentError: You must provide a url to stream and a message to show.")
		}
		// Remove the angled brackets <> surrounding a URL that prevent an embed from appearing
		CleanURL := strings.Replace(strings.Replace(UrlOrMessage[0], "<", "", 1), ">", "", 1)
		ctx.Session.UpdateStreamingStatus(0, strings.Join(UrlOrMessage[1:], " "), CleanURL)
		_, err = ctx.Send("Set `stream` status to `" + strings.Join(UrlOrMessage[1:], " ") + "` at URL <" + CleanURL + ">")
	case "reset":
		ctx.Session.UpdateStatus(0, "Go, Discord, go!")
		_, err = ctx.Send("Presence reset.")
	default:
		_, err = ctx.Send(ActivityType + " is not an option. Select from `play`, `listen`, `stream`, `reset`")
	}
	return
}

func init() {
	cog := commands.NewCog("Owner", "Developer restricted commands", true)
	cog.AddCommand("presence", "Changes the bot's presence", "<ActivityType> [url:stream] [message]", presence).
		SetDefaultArg([]string{})
	cog.Load()
}
