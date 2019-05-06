package commands

import discord "github.com/bwmarrin/discordgo"

// Context command context
type Context struct {
	Session *discord.Session
	Message *discord.MessageCreate
	Guild   *discord.Guild
	Channel *discord.Channel
	Author  *discord.User
	Args    []string
	Prefix  string
}
