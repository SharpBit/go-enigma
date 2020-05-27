package cogs

import (
	"strings"

	"github.com/SharpBit/go-enigma/commands"
	"github.com/SharpBit/go-enigma/utils"
	"github.com/bwmarrin/discordgo"
)

func ban(ctx *commands.Context, member *discordgo.Member, reason ...string) (err error) {
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

func kick(ctx *commands.Context, member *discordgo.Member, reason ...string) (err error) {
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
		AddCheck(utils.PermCheck("Ban Members", discordgo.PermissionBanMembers)).
		AddCheck(utils.BotPermCheck("Ban Members", discordgo.PermissionBanMembers))
	cog.AddCommand("unban", "Unban a user from the guild", "<NameOrID> [reason]", unban).
		SetDefaultArg("None").
		AddCheck(utils.PermCheck("Ban Members", discordgo.PermissionBanMembers)).
		AddCheck(utils.BotPermCheck("Ban Members", discordgo.PermissionBanMembers))
	cog.AddCommand("kick", "Kick a member from the guild", "<member> [reason]", kick).
		SetDefaultArg("None").
		AddCheck(utils.PermCheck("Kick Members", discordgo.PermissionKickMembers)).
		AddCheck(utils.BotPermCheck("Kick Members", discordgo.PermissionKickMembers))
	cog.Load()

}
