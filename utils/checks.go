package utils

import (
	"fmt"

	"github.com/SharpBit/go-enigma/commands"
	discord "github.com/bwmarrin/discordgo"
)

func PermCheck(name string, RequiredPerm int) commands.CheckFunction {
	return func(ctx *commands.Context) (bool, error) {
		state := ctx.Session.State
		perms, err := state.UserChannelPermissions(ctx.Author.ID, ctx.Channel.ID)
		if err != nil {
			return false, err
		}

		if perms&RequiredPerm == RequiredPerm {
			return true, nil
		}
		return false, fmt.Errorf("CommandCheckError: PermissionError: You need the **" + name + "** permission to perform this command.")
	}
}

func BotPermCheck(name string, RequiredPerm int) commands.CheckFunction {
	return func(ctx *commands.Context) (bool, error) {
		state := ctx.Session.State
		perms, err := state.UserChannelPermissions(ctx.Session.State.User.ID, ctx.Channel.ID)
		if err != nil {
			return false, err
		}

		if perms&RequiredPerm == RequiredPerm {
			return true, nil
		}
		return false, fmt.Errorf("CommandCheckError: BotPermissionError: I need the **" + name + "** permission to perform this command.")
	}
}

// OwnerOnly: Check to see if the author is the bot owner
func OwnerOnly(ctx *commands.Context) (bool, error) {
	OwnerID := GetConfig("ownerID")
	if ctx.Author.ID != OwnerID {
		return false, fmt.Errorf("CommandCheckError: You are not the owner of this bot!")
	}
	return true, nil
}

func GuildOnly(ctx *commands.Context) (bool, error) {
	if ctx.Channel.Type == discord.ChannelTypeDM || ctx.Channel.Type == discord.ChannelTypeGroupDM {
		return false, fmt.Errorf("CommandCheckError: Command must be run in a guild.")
	}
	return true, nil
}
