package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Context: A class that stores information about the message and the Command
type Context struct {
	Session *discordgo.Session
	Message *discordgo.MessageCreate
	Guild   *discordgo.Guild
	Channel *discordgo.Channel
	Author  *discordgo.User
	Prefix  string
	Command *Command
}

// Send a message to ctx.Channel
func (ctx *Context) Send(content string) (*discordgo.Message, error) {
	return ctx.Session.ChannelMessageSend(ctx.Channel.ID, content)
}

// SendComplex: Send an embed/complex message to ctx.Channel
func (ctx *Context) SendComplex(content string, embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	data := &discordgo.MessageSend{Content: content, Embed: embed}
	return ctx.Session.ChannelMessageSendComplex(ctx.Channel.ID, data)
}

// SendError replies with the error and help if sendHelp is true
func (ctx *Context) SendError(err error, sendHelp bool) (*discordgo.Message, error) {
	xmark, err2 := ctx.GetEmoji("xmark")
	if err2 != nil {
		return nil, err2
	}

	if sendHelp {
		usageString := "`" + ctx.Prefix + ctx.Command.Name
		if ctx.Command.Usage != "" {
			usageString += " " + ctx.Command.Usage
		}
		usageString += "`"
		em := NewEmbed().
			SetColor(0xe74c3c).
			SetTitle(usageString).
			SetDescription(ctx.Command.Description).
			MessageEmbed
		return ctx.SendComplex(xmark.MessageFormat()+" "+err.Error(), em)
	}
	return ctx.Send(xmark.MessageFormat() + " " + err.Error())
}

// Edit a message content
func (ctx *Context) Edit(msg *discordgo.Message, content string) *discordgo.Message {
	edit := discordgo.NewMessageEdit(msg.ChannelID, msg.ID)
	edit.SetContent(content)
	return msg
}

// EditComplex: Edit a message's content an embed
func (ctx *Context) EditComplex(msg *discordgo.Message, content string, embed *discordgo.MessageEmbed) *discordgo.Message {
	edit := discordgo.NewMessageEdit(msg.ChannelID, msg.ID)
	edit.SetContent(content).SetEmbed(embed)
	return msg
}

// CodeBlock returns code formatted into a codeblock to send to Discord
func (ctx *Context) CodeBlock(content string, lang string) (formatted string) {
	return "```" + lang + "\n" + content + "\n```"
}

// GetBan: Checks the bans of ctx.Guild and returns a string (User ID)
func (ctx *Context) GetBan(input string) (userID string, err error) {

	bans, err := ctx.Session.GuildBans(ctx.Guild.ID)
	if err != nil {
		return "", fmt.Errorf("BotPermissionError: Do not have ban members permissions.")
	}

	for _, b := range bans {
		if len(input) > 5 && b.User.Username == input[:len(input)-5] && b.User.Discriminator == input[len(input)-4:] {
			return b.User.ID, nil
		}

		if b.User.ID == input {
			return b.User.ID, nil
		}
	}
	return "", fmt.Errorf("NotFoundError: no ban found")
}

func (ctx *Context) GetEmoji(name string) (emoji *discordgo.Emoji, err error) {
	GuildID := "571500500357480448"
	guild, err := ctx.Session.State.Guild(GuildID)
	if err != nil {
		guild, err = ctx.Session.Guild(GuildID)
		if err != nil {
			return nil, err
		}
	}

	for _, e := range guild.Emojis {
		CleanedName := strings.Split(e.Name, " ")[0]
		if CleanedName == name {
			return e, nil
		}
	}

	return nil, fmt.Errorf("error ctx.GetEmoji: Emoji not found.")
}
