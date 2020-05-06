package events

import (
	commands "github.com/SharpBit/go-enigma/cogs"

	discord "github.com/bwmarrin/discordgo"
)

// MessageCreate : Event that's called whenever a message is sent
func MessageCreate(session *discord.Session, msg *discord.MessageCreate) {
	ctx, err := commands.HandleCommands(session, msg)
	if err != nil {
		commands.HandleCommandError(ctx, err)
	}
}
