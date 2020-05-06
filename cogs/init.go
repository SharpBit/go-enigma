package cogs

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/SharpBit/go-enigma/utils"

	discord "github.com/bwmarrin/discordgo"
)

var (
	prefix  = utils.GetConfig("prefix")
	OwnerID = utils.GetConfig("owner_id")
)

// Context: A class that stores information about the message and the Command
type Context struct {
	Session *discord.Session
	Message *discord.MessageCreate
	Guild   *discord.Guild
	Channel *discord.Channel
	Author  *discord.User
	Prefix  string
	Command *Command
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

// SendError replies with the error and help if sendHelp is true
func (ctx *Context) SendError(err error, sendHelp bool) (*discord.Message, error) {
	usageString := "`" + ctx.Prefix + ctx.Command.Name + " " + ctx.Command.Usage + "`"
	em := utils.NewEmbed().
		SetColor(0xe74c3c).
		SetTitle(usageString).
		SetDescription(ctx.Command.Description).
		MessageEmbed
	data := &discord.MessageSend{Content: err.Error(), Embed: em}
	return ctx.Session.ChannelMessageSendComplex(ctx.Channel.ID, data)
}

// CodeBlock returns code formatted into a codeblock to send to Discord
func (ctx *Context) CodeBlock(content string, lang string) (formatted string) {
	return "```" + lang + "\n" + content + "\n```"
}

// GetBan: Checks the guild's bans and returns a string (User ID)
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
	Name           string
	Description    string
	Aliases        []string
	Dev            bool
	Usage          string
	Run            reflect.Value
	HasOptionalArg bool
	DefaultArg     reflect.Value
}

func (cmd *Command) Type() reflect.Type {
	return reflect.TypeOf(cmd.Run)
}

func (cmd *Command) SetUsage(usage string) *Command {
	cmd.Usage = usage
	return cmd
}

func (cmd *Command) SetAliases(aliases ...string) *Command {
	for _, alias := range aliases {
		// Check to see if the alias is already registered as a command or alias
		_, existing := AliasMap[alias]
		if !existing {
			_, existing = CommandMap[alias]
			if !existing {
				// Not registered yet, add it to the aliases
				cmd.Aliases = append(cmd.Aliases, alias)
			}
		}
	}
	return cmd
}

// SetDefaultArg: Makes the last argument optional and will pass in a default argument if not provided
func (cmd *Command) SetDefaultArg(def interface{}) *Command {
	cmd.DefaultArg = reflect.ValueOf(def)
	cmd.HasOptionalArg = true
	return cmd
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
func (cog *Cog) AddCommand(name, description string, usage string, run interface{}) *Command {
	cmd, existing := NewCommand(name, description)
	if existing {
		fmt.Println("error: command " + name + " already exists")
	}
	cmd.Run = reflect.ValueOf(run)
	if cog.Dev == true {
		cmd.Dev = true
	}
	cmd.SetUsage(usage)
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

// HandleCommands gets called on messageCreate
func HandleCommands(session *discord.Session, msg *discord.MessageCreate) (ctx *Context, err error) {
	// so the bot doesn't respond to other bots or webhooks (including itself)
	if msg.Author.Bot || msg.WebhookID != "" {
		return nil, nil
	}

	// Check if the message starts with the bot prefix
	if !strings.HasPrefix(msg.Content, prefix) {
		return nil, nil
	}

	// gets the message's guild
	guild, err := session.State.Guild(msg.GuildID)
	if err != nil {
		guild, err = session.Guild(msg.GuildID)
		if err != nil {
			return nil, err
		}
	}

	// gets the message's channel
	channel, err := session.State.Channel(msg.ChannelID)
	if err != nil {
		channel, err = session.Channel(msg.ChannelID)
		if err != nil {
			return nil, err
		}
	}

	ctx = &Context{
		Session: session,
		Message: msg,
		Guild:   guild,
		Channel: channel,
		Author:  msg.Author,
		Prefix:  prefix,
	}

	// Separates the commands from the arguments
	input := strings.Fields(msg.Content)
	CmdString, args := strings.Trim(input[0], prefix), input[1:]

	cmd, ok := CommandMap[CmdString]
	if !ok {
		cmdName, ok := AliasMap[CmdString]
		if !ok {
			return nil, fmt.Errorf("InvokeCommandError: Invalid command name: %s", cmdName)
		}

		// Use the command name retrieved from the AliasMap to get the Command
		cmd = CommandMap[cmdName]
	}

	ctx.Command = cmd

	// Ignore developer commands if the owner did not invoke the command
	if cmd.Dev == true && msg.Author.ID != OwnerID {
		return ctx, fmt.Errorf("PermissionError: Author is not owner")
	}

	// Allow multiple word arguments as long as they are surrounded by quotes.
	// i.e.  !ban "Sharp Bit#0001" spamming  would return the arguments []string{"Sharp Bit#0001", "spamming"}
	var ParsedArgs []string
	var currentParsed string
	for _, arg := range args {
		// Arg starts with "
		// Add all future args to currentParsed until another " is found
		if strings.HasPrefix(arg, "\"") && currentParsed == "" {
			// The word with quote removed followed by a space
			currentParsed += arg[1:] + " "
		} else if strings.HasSuffix(arg, "\"") && currentParsed != "" {
			// The end of the argument is found, add the word without the quote and add it to ParsedArgs
			currentParsed += arg[:len(arg)-1]
			ParsedArgs = append(ParsedArgs, currentParsed)
			currentParsed = ""
		} else if currentParsed != "" {
			// Add the word to currentParsed, the end is not found
			currentParsed += arg + " "
		} else {
			// Not a multi-word arg: add the arg
			ParsedArgs = append(ParsedArgs, arg)
		}
	}

	cmdType := cmd.Run.Type()

	// We have to check argument numbers before converting argument types so it doesn't try to convert extra arguments
	trueLength := len(ParsedArgs) + 1 // to account for ctx
	if cmdType.IsVariadic() {
		// Variadic arguments
		// Last argument is optional
		if cmd.HasOptionalArg {
			// Less than the minimum number of parameters
			if trueLength < cmdType.NumIn()-1 {
				return ctx, fmt.Errorf("ArgumentError: Incorrect number of arguments")
			}
		} else {
			// Not equal number of args
			if trueLength < cmdType.NumIn() {
				return ctx, fmt.Errorf("ArgumentError: Incorrect number of arguments")
			}
		}
	} else {
		// Set number of arguments
		// Last argument is optional
		if cmd.HasOptionalArg {
			// Not equal to or one less
			if !(trueLength == cmdType.NumIn() || trueLength == cmdType.NumIn()-1) {
				return ctx, fmt.Errorf("ArgumentError: Incorrect number of arguments")
			}
		} else {
			// Not equal number of args
			if trueLength != cmdType.NumIn() {
				return ctx, fmt.Errorf("ArgumentError: Incorrect number of arguments")
			}
		}
	}

	ConvertedArgs := []reflect.Value{reflect.ValueOf(ctx)}

	for i, arg := range ParsedArgs {
		// Skip the context argument, decremented later
		i++

		t := cmdType.In(i).Kind()

		if i >= cmdType.NumIn() {
			// Use the last argument's value since it is a slice
			i = cmdType.NumIn() - 1
		}

		// If it is a slice, use the type of the slice
		if t == reflect.Slice {
			t = cmdType.In(i).Elem().Kind()
		}

		fmt.Println(t)

		switch t {
		// String
		case reflect.String:
			ConvertedArgs = append(ConvertedArgs, reflect.ValueOf(arg))
		// Int
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val, err := strconv.ParseInt(arg, 10, 64)
			if err != nil {
				return ctx, err
			}
			ConvertedArgs = append(ConvertedArgs, reflect.ValueOf(val))
		// Unsigned int (only positive)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val, err := strconv.ParseUint(arg, 10, 64)
			if err != nil {
				return ctx, err
			}
			ConvertedArgs = append(ConvertedArgs, reflect.ValueOf(val))
		// Floats
		case reflect.Float32, reflect.Float64:
			val, err := strconv.ParseFloat(arg, 64)
			if err != nil {
				return ctx, err
			}
			ConvertedArgs = append(ConvertedArgs, reflect.ValueOf(val))
		// Parse common boolean indicators
		case reflect.Bool:
			switch strings.ToLower(arg) {
			case "true", "t", "yes", "y", "1":
				ConvertedArgs = append(ConvertedArgs, reflect.ValueOf(true))
			case "false", "f", "no", "n", "0":
				ConvertedArgs = append(ConvertedArgs, reflect.ValueOf(false))
			default:
				return ctx, fmt.Errorf("ArgumentError: Invalid boolean.")
			}
		// It is a pointer, get the type
		case reflect.Ptr:
			fmt.Println(cmdType.In(i).Elem().Name())
			if cmdType.In(i).Elem().Name() == "User" {
				// *discord.User
				var user *discord.User

				// Look for mentions
				if len(msg.Mentions) > 0 {
					user = msg.Mentions[0]
				}

				for _, m := range guild.Members {
					// Look for Name#Discrim
					if arg == m.User.String() {
						user, err = session.User(m.User.ID)
						if err != nil {
							return ctx, err
						}
					}

					// Look for ID
					if arg == m.User.ID {
						user, err = session.User(m.User.ID)
						if err != nil {
							return ctx, err
						}
					}
				}

				if user != nil {
					ConvertedArgs = append(ConvertedArgs, reflect.ValueOf(user))
					break
				}

				return ctx, fmt.Errorf("ArgumentError: User was not found.")
			} else if cmdType.In(i).Elem().Name() == "Member" {
				// *discord.Member
				var member *discord.Member

				// Look for mentions
				if len(msg.Mentions) > 0 {
					arg = msg.Mentions[0].ID
				}

				for _, m := range guild.Members {
					// Look for Name#Discrim
					if arg == m.User.String() {
						member = m
					}

					// Look for ID
					if arg == m.User.ID {
						member = m
					}
				}

				if member != nil {
					ConvertedArgs = append(ConvertedArgs, reflect.ValueOf(member))
					break
				}

				return ctx, fmt.Errorf("ArgumentError: Member was not found.")
			} else {
				// Not an implemented pointer type
				return ctx, fmt.Errorf("Invalid pointer type")
			}
		default:
			// Not an implemented type
			return ctx, fmt.Errorf("ArgmentError: Invalid type")
		}
		i--
	}

	// Add the default arg if there is not already a value
	if len(ConvertedArgs) < cmdType.NumIn() && cmd.HasOptionalArg {
		if cmd.DefaultArg == reflect.ValueOf(nil) {
			t := cmdType.In(cmdType.NumIn() - 1)
			ConvertedArgs = append(ConvertedArgs, reflect.New(t).Elem())
		} else {
			ConvertedArgs = append(ConvertedArgs, cmd.DefaultArg)
		}
	}

	// Call the function and return any errors to the event that called the function (messageCreate)
	// so it can call the Error Handler
	res := cmd.Run.Call(ConvertedArgs)
	if res[0].Interface() != nil {
		err = res[0].Interface().(error)
	} else {
		err = nil
	}
	return ctx, err
}

// HandleCommandError: Handles and recovers from errors without panicking
func HandleCommandError(ctx *Context, err error) {
	// No command found
	if strings.HasPrefix(err.Error(), "InvokeCommandError") {
		return
	}

	// Guild or channel could not be retrieved
	if ctx == nil {
		panic(err)
	}

	// Something is wrong with the user inputted arguments
	if strings.HasPrefix(err.Error(), "ArgumentError") {
		ctx.SendError(err, true)
		return
	}

	// Log any other errors without panicking
	fmt.Println(err)
	debug.PrintStack()
}
