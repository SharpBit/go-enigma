package cogs

import (
	discord "github.com/bwmarrin/discordgo"
	"strings"
)

func ban(ctx *Context, member *discord.Member, reason ...string) (err error) {
	ReasonFmt := ctx.Author.Username + "#" + ctx.Author.Discriminator + " (" + strings.Join(reason, " ") + ")"
	ctx.Session.GuildBanCreateWithReason(ctx.Guild.ID, member.User.ID, ReasonFmt, 1)
	_, err = ctx.Send("Done.")
	return
}

func unban(ctx *Context, NameOrID string, reason ...string) (err error) {
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

func kick(ctx *Context, member *discord.Member, reason ...string) (err error) {
	ReasonFmt := ctx.Author.Username + "#" + ctx.Author.Discriminator + " (" + strings.Join(reason, " ") + ")"
	err = ctx.Session.GuildMemberDeleteWithReason(ctx.Guild.ID, member.User.ID, ReasonFmt)
	if err != nil {
		return
	}
	_, err = ctx.Send("Done.")
	return
}

func init() {
	cog := NewCog("Mod", "Guild Moderation commands", false)
	cog.AddCommand("ban", "Ban a member from the guild", "<member> [reason]", ban).
		SetDefaultArg("None")
	cog.AddCommand("unban", "Unban a user from the guild", "<NameOrID> [reason]", unban).
		SetDefaultArg("None")
	cog.AddCommand("kick", "Kick a member from the guild", "<member> [reason]", kick).
		SetDefaultArg("None")
	cog.Load()

}
