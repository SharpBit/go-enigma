package cogs

import (
	"strings"

	"github.com/SharpBit/go-enigma/commands"
	"github.com/SharpBit/go-enigma/utils"
	discord "github.com/bwmarrin/discordgo"
)

func ban(ctx *commands.Context, member *discord.Member, reason ...string) (err error) {
	ReasonFmt := ctx.Author.Username + "#" + ctx.Author.Discriminator + " (" + strings.Join(reason, " ") + ")"
	ctx.Session.GuildBanCreateWithReason(ctx.Guild.ID, member.User.ID, ReasonFmt, 1)
	_, err = ctx.Send("Done.")
	return
}

func unban(ctx *commands.Context, NameOrID string, reason ...string) (err error) {
	user, err := ctx.GetBan(NameOrID)
	if err != nil {
		ctx.SendError(err, false)
		return
	}

	err = ctx.Session.GuildBanDelete(ctx.Guild.ID, user)
	if err != nil {
		ctx.SendError(err, false)
		return
	}
	_, err = ctx.Send("Done.")
	return
}

func kick(ctx *commands.Context, member *discord.Member, reason ...string) (err error) {
	ReasonFmt := ctx.Author.Username + "#" + ctx.Author.Discriminator + " (" + strings.Join(reason, " ") + ")"
	err = ctx.Session.GuildMemberDeleteWithReason(ctx.Guild.ID, member.User.ID, ReasonFmt)
	if err != nil {
		return
	}
	_, err = ctx.Send("Done.")
	return
}

func init() {
	cog := commands.NewCog("Mod", "Guild Moderation commands")
	cog.AddCommand("ban", "Ban a member from the guild", "<member> [reason]", ban).
		SetDefaultArg("None").
		AddCheck(utils.PermCheck(discord.PermissionBanMembers)).
		AddCheck(utils.BotPermCheck(discord.PermissionBanMembers))
	cog.AddCommand("unban", "Unban a user from the guild", "<NameOrID> [reason]", unban).
		SetDefaultArg("None").
		AddCheck(utils.PermCheck(discord.PermissionBanMembers)).
		AddCheck(utils.BotPermCheck(discord.PermissionBanMembers))
	cog.AddCommand("kick", "Kick a member from the guild", "<member> [reason]", kick).
		SetDefaultArg("None").
		AddCheck(utils.PermCheck(discord.PermissionKickMembers)).
		AddCheck(utils.BotPermCheck(discord.PermissionKickMembers))
	cog.Load()

}
