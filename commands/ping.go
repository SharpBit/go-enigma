package commands

import (
	discord "github.com/bwmarrin/discordgo"
)

// Ping command
func Ping(session *discord.Session, msg *discord.Message) {
	session.ChannelMessageSend(msg.ChannelID, "Pong!")
}
