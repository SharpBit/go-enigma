package cogs

import (
	"errors"
	"fmt"

	discord "github.com/bwmarrin/discordgo"
)

// Context : command context
type Context struct {
	Session *discord.Session
	Message *discord.MessageCreate
	Guild   *discord.Guild
	Channel *discord.Channel
	Author  *discord.User
	Args    []string
	Prefix  string
}

// Send a message to the channel
func (ctx *Context) Send(content string) (*discord.Message, error) {
	return ctx.Session.ChannelMessageSend(ctx.Channel.ID, content)
}

// SendComplex an embed/complex message to the channel
func (ctx *Context) SendComplex(content string, embed *discord.MessageEmbed) (*discord.Message, error) {
	data := &discord.MessageSend{Content: content, Embed: embed}
	return ctx.Session.ChannelMessageSendComplex(ctx.Channel.ID, data)
}

// CodeBlock returns code formatted into a codeblock to send to Discord
func (ctx *Context) CodeBlock(content string, lang string) (formatted string) {
	return "```" + lang + "\n" + content + "\n```"
}

// ParseUser gets User object from a Mention, Name#Discrim, or ID
func (ctx *Context) ParseUser(input string) (user *discord.User, err error) {
	if len(ctx.Message.Mentions) > 0 {
		return ctx.Message.Mentions[0], nil
	}

	for _, m := range ctx.Guild.Members {
		if m.User.Username == input[:len(input)-5] && m.User.Discriminator == input[len(input)-4:] {
			return ctx.Session.User(m.User.ID)
		}

		if m.User.ID == input {
			return ctx.Session.User(m.User.ID)
		}
	}
	return nil, errors.New("error ParseUser: no user found")
}

// GetBan does the same as ParseUser but checks the guild's bans and returns a string (User ID)
func (ctx *Context) GetBan(input string) (userID string, err error) {

	bans, err := ctx.Session.GuildBans(ctx.Guild.ID)
	if err != nil {
		ctx.Send("Do not have ban members permissions")
		return
	}

	for _, b := range bans {
		if b.User.Username == input[:len(input)-5] && b.User.Discriminator == input[len(input)-4:] {
			return b.User.ID, nil
		}

		if b.User.ID == input {
			return b.User.ID, nil
		}
	}
	return "", errors.New("error GetBan: no ban found")
}

/*
Command Structs and Functions
*/

// CommandMap is a map that gets the user's command input and retrieves its respective function
var CommandMap = make(map[string]*Command)

// AliasMap finds the command of each alias
var AliasMap = make(map[string]string)

// CogMap finds the Cog object from the name
var CogMap = make(map[string]*Cog)

// Command is a command object
type Command struct {
	Name        string
	Description string
	Aliases     []string
	Dev         bool
	Run         func(*Context)
}

// NewCommand creates a new command
func NewCommand(name, description string) (cmd *Command, existing bool) {
	_, existing = CommandMap[name]
	if existing {
		return nil, existing
	}
	cmd = &Command{Name: name, Description: description}
	return cmd, existing
}

// RegisterCommand adds the command to the CommandMap
func RegisterCommand(cmd *Command) {
	CommandMap[cmd.Name] = cmd
	for _, alias := range cmd.Aliases {
		AliasMap[alias] = cmd.Name
	}
}

// UnregisterCommand removes the command from the CommandMap
func UnregisterCommand(cmd *Command) {
	delete(CommandMap, cmd.Name)
	for _, alias := range cmd.Aliases {
		delete(AliasMap, alias)
	}
}

/*
Cog structs and functions
*/

// Cog is a similar class to commands.Cog in discord.py
type Cog struct {
	Name        string
	Description string
	Dev         bool
	Commands    []*Command
	Loaded      bool
}

// NewCog creates a new cog instance
func NewCog(name, description string, dev bool) *Cog {
	return &Cog{Name: name, Description: description, Dev: dev, Loaded: false}
}

// AddCommand : Adds a command to the cog
func (cog *Cog) AddCommand(name, description string, aliases []string, run func(*Context)) *Command {
	cmd, existing := NewCommand(name, description)
	if existing {
		fmt.Println("error: command " + name + " already exists")
	}
	cmd.Run = run
	cmd.Aliases = aliases
	if cog.Dev == true {
		cmd.Dev = true
	}
	cog.Commands = append(cog.Commands, cmd)

	return cmd
}

// Load : Registers each command in the cog
func (cog *Cog) Load() {
	for _, cmd := range cog.Commands {
		RegisterCommand(cmd)
	}
	CogMap[cog.Name] = cog
	cog.Loaded = true
}

// Unload : Unregisters each command in the cog
func (cog *Cog) Unload() {
	for _, cmd := range cog.Commands {
		UnregisterCommand(cmd)
	}
	delete(CogMap, cog.Name)
	cog.Loaded = false
}
