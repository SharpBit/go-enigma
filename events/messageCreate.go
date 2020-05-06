package events

import (
	"github.com/SharpBit/go-enigma/handlers"

	discord "github.com/bwmarrin/discordgo"
)

// MessageCreate : Event that's called whenever a message is sent
func MessageCreate(session *discord.Session, msg *discord.MessageCreate) {
	ctx, err := handlers.HandleCommands(session, msg)
	if err != nil {
		handlers.HandleCommandError(ctx, err)
	}
}
