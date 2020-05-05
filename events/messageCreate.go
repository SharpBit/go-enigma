package events

import (
	"strings"

	commands "github.com/SharpBit/go-enigma/cogs"

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
	input := strings.Fields(msg.Content)
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

	// Allow multiple word arguments as long as they are surrounded by quotes.
	// i.e.  !ban Abigail Brown#0001 spamming  would return the arguments ["Abigail Brown#0001", "spamming"]
	var ParsedArgs []string
	var currentParsed string
	for _, arg := range args {
		if strings.HasPrefix(arg, "\"") && currentParsed == "" {
			currentParsed += arg[1:] + " "
		} else if strings.HasSuffix(arg, "\"") && currentParsed != "" {
			currentParsed += arg[:len(arg)-1]
			ParsedArgs = append(ParsedArgs, currentParsed)
			currentParsed = ""
		} else if currentParsed != "" {
			currentParsed += arg + " "
		} else {
			ParsedArgs = append(ParsedArgs, arg)
		}
	}

	cmd.Run(&commands.Context{
		Session: session,
		Message: msg,
		Guild:   guild,
		Channel: channel,
		Author:  msg.Author,
		Args:    ParsedArgs,
		Prefix:  prefix,
	})
}
