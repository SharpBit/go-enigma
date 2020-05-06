package utils

import (
	"fmt"

	"github.com/SharpBit/go-enigma/commands"
	discord "github.com/bwmarrin/discordgo"
)

// OwnerOnly: Check to see if
func OwnerOnly(ctx *commands.Context) (bool, error) {
	OwnerID := GetConfig("ownerID")
	if ctx.Author.ID != OwnerID {
		return false, fmt.Errorf("CommandCheckError: You are not the owner of this bot")
	}
	return true, nil
}

func GuildOnly(ctx *commands.Context) (bool, error) {
	if ctx.Channel.Type == discord.ChannelTypeDM || ctx.Channel.Type == discord.ChannelTypeGroupDM {
		return false, fmt.Errorf("CommandCheckError: Command must be run in a guild.")
	}
	return true, nil
}
