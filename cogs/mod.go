package cogs

import "strings"

func ban(ctx *Context) {
	if len(ctx.Args) == 0 {
		ctx.Send("Please pass in a member to ban.")
		return
	}
	member, err := ctx.ParseUser(ctx.Args[0])
	if err != nil {
		ctx.Send("Member not found.")
		return
	}

	reason := ctx.Author.Username + "#" + ctx.Author.Discriminator
	if len(ctx.Args) > 1 {
		reason = reason + " (" + strings.Join(ctx.Args[1:], " ") + ")"
	} else {
		reason = reason + " (None)"
	}
	ctx.Session.GuildBanCreateWithReason(ctx.Guild.ID, member.ID, reason, 1)
	ctx.Send("Done.")
}

func unban(ctx *Context) {
	if len(ctx.Args) == 0 {
		ctx.Send("Please pass in a member to unban.")
		return
	}

	user, err := ctx.GetBan(ctx.Args[0])
	if err != nil {
		ctx.Send("No ban entry " + ctx.Args[0] + " found.")
		return
	}

	ctx.Session.GuildBanDelete(ctx.Guild.ID, user)
	ctx.Send("Done.")
}

func kick(ctx *Context) {
	if len(ctx.Args) == 0 {
		ctx.Send("Please pass in a member to kick.")
		return
	}
	member, err := ctx.ParseUser(ctx.Args[0])
	if err != nil {
		ctx.Send("Member not found.")
		return
	}

	reason := ctx.Author.Username + "#" + ctx.Author.Discriminator
	if len(ctx.Args) > 1 {
		reason = reason + " (" + strings.Join(ctx.Args[1:], " ") + ")"
	}
	ctx.Session.GuildMemberDeleteWithReason(ctx.Guild.ID, member.ID, reason)
	ctx.Send("Done.")
}

func init() {
	cog := NewCog("Mod", "Guild Moderation commands", false)
	cog.AddCommand("ban", "Ban a member from the guild", []string{}, ban)
	cog.AddCommand("unban", "Unban a user from the guild", []string{}, unban)
	cog.AddCommand("kick", "Kick a member from the guild", []string{}, kick)
	cog.Load()

}
