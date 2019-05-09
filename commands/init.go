package commands

import (
	"fmt"

	discord "github.com/bwmarrin/discordgo"
)

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

/*
Command Structs and Functions
*/

// CommandMap is a map that gets the user's command input and retrieves its respective function
var CommandMap = make(map[string]*Command)

// AliasMap finds the commands of each alias
var AliasMap = make(map[string]string)

// Command is a command object
type Command struct {
	Name        string
	Description string
	Aliases     []string
	Dev         bool
	Run         func(*Context)
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
func UnregisterCommand(cmd string) {
	delete(CommandMap, cmd)
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
}

// NewCog creates a new cog instance
func NewCog(name, description string, dev bool) *Cog {
	return &Cog{Name: name, Description: description, Dev: dev}
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
}

// Unload : Unregisters each command in the cog
func (cog *Cog) Unload() {
	for _, cmd := range cog.Commands {
		UnregisterCommand(cmd.Name)
	}
}
