package handlers

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/SharpBit/go-enigma/commands"
	"github.com/SharpBit/go-enigma/utils"
	"github.com/bwmarrin/discordgo"
)

var (
	prefix  = utils.GetConfig("prefix")
	OwnerID = utils.GetConfig("ownerID")
)

// HandleCommands gets called on messageCreate
func HandleCommands(session *discordgo.Session, msg *discordgo.MessageCreate) (ctx *commands.Context, err error) {
	// so the bot doesn't respond to other bots or webhooks (including itself)
	if msg.Author.Bot || msg.WebhookID != "" {
		return nil, nil
	}

	// Check if the message starts with the bot prefix
	if !strings.HasPrefix(msg.Content, prefix) {
		return nil, nil
	}

	// gets the message's channel
	channel, err := session.State.Channel(msg.ChannelID)
	if err != nil {
		channel, err = session.Channel(msg.ChannelID)
		if err != nil {
			return nil, err
		}
	}

	ctx = &commands.Context{
		Session: session,
		Message: msg,
		Channel: channel,
		Author:  msg.Author,
		Prefix:  prefix,
	}

	// gets the message's guild
	guild, err := session.State.Guild(msg.GuildID)
	if err != nil {
		guild, err = session.Guild(msg.GuildID)
		if err != nil {
			if !(channel.Type == discordgo.ChannelTypeDM || channel.Type == discordgo.ChannelTypeGroupDM) {
				return nil, err
			}
		}
	}

	ctx.Guild = guild

	// Separates the commands from the arguments
	input := strings.Fields(msg.Content)
	CmdString, args := strings.Trim(input[0], prefix), input[1:]

	cmd, ok := commands.CommandMap[CmdString]
	if !ok {
		cmdName, ok := commands.AliasMap[CmdString]
		if !ok {
			return nil, fmt.Errorf("InvokeCommandError: Invalid command name: %s", cmdName)
		}

		// Use the command name retrieved from the AliasMap to get the Command
		cmd = commands.CommandMap[cmdName]
	}

	ctx.Command = cmd

	for _, check := range cmd.Checks {
		passed, err := check(ctx)
		if !passed {
			return ctx, err
		}
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
				return ctx, fmt.Errorf("ArgumentError: Incorrect number of arguments provided.")
			}
		} else {
			// Not equal number of args
			if trueLength < cmdType.NumIn() {
				return ctx, fmt.Errorf("ArgumentError: Incorrect number of arguments provided.")
			}
		}
	} else {
		// Set number of arguments
		// Last argument is optional
		if cmd.HasOptionalArg {
			// Not equal to or one less
			if !(trueLength == cmdType.NumIn() || trueLength == cmdType.NumIn()-1) {
				return ctx, fmt.Errorf("ArgumentError: Incorrect number of arguments provided.")
			}
		} else {
			// Not equal number of args
			if trueLength != cmdType.NumIn() {
				return ctx, fmt.Errorf("ArgumentError: Incorrect number of arguments provided.")
			}
		}
	}

	ConvertedArgs := []reflect.Value{reflect.ValueOf(ctx)}

	for i, arg := range ParsedArgs {
		// Skip the context argument, decremented later
		i++

		if i >= cmdType.NumIn() {
			// Use the last argument's value since it is a slice
			i = cmdType.NumIn() - 1
		}

		t := cmdType.In(i).Kind()

		// If it is a slice, use the type of the slice
		if t == reflect.Slice {
			t = cmdType.In(i).Elem().Kind()
		}

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
				return ctx, fmt.Errorf("ArgumentError: Invalid boolean value.")
			}
		// It is a pointer, get the type
		case reflect.Ptr:
			if cmdType.In(i).Elem().Name() == "User" {
				// *discord.User
				var user *discordgo.User

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
				var member *discordgo.Member

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
				return ctx, fmt.Errorf("ArgumentError: Invalid pointer type.")
			}
		default:
			// Not an implemented type
			return ctx, fmt.Errorf("ArgmentError: Invalid type.")
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
