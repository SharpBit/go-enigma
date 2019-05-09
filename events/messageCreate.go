package events

import (
	"strings"

	commands "../cogs"

	discord "github.com/bwmarrin/discordgo"
)

var (
	prefix = "//"
	owner  = "281821029490229251"
)

// MessageCreate event
func MessageCreate(session *discord.Session, msg *discord.MessageCreate) {

	// so the bot doesn't respond to other bots or webhooks (including itself)
	if msg.Author.Bot || msg.WebhookID != "" {
		return
	}

	if !strings.HasPrefix(msg.Content, prefix) {
		return
	}

	// gets the message's guild
	guild, err := session.State.Guild(msg.GuildID)
	if err != nil {
		guild, err = session.Guild(msg.GuildID)
		if err != nil {
			return
		}
	}

	// gets the message's channel
	channel, err := session.State.Channel(msg.ChannelID)
	if err != nil {
		channel, err = session.Channel(msg.ChannelID)
		if err != nil {
			return
		}
	}

	// Separates the commands from the arguments
	input := strings.Split(msg.Content, " ")
	CmdString, args := strings.Trim(input[0], prefix), input[1:]

	cmd, ok := commands.CommandMap[CmdString]
	if !ok {
		alias, ok := commands.AliasMap[CmdString]
		if !ok {
			return
		}

		cmd = commands.CommandMap[alias]
	}

	if cmd.Dev == true && msg.Author.ID != "281821029490229251" {
		return
	}

	cmd.Run(&commands.Context{
		Session: session,
		Message: msg,
		Guild:   guild,
		Channel: channel,
		Author:  msg.Author,
		Args:    args,
		Prefix:  prefix,
	})
}
