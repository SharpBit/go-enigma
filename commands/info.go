package commands

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"../utils"

	discord "github.com/bwmarrin/discordgo"
)

func serverinfo(ctx *Context) {
	guild := ctx.Guild
	channel := ctx.Channel

	if channel.Type == discord.ChannelTypeDM || channel.Type == discord.ChannelTypeGroupDM {
		return
	}

	onlineMembers := 0
	for _, m := range guild.Presences {
		if m.Status != discord.StatusOffline {
			onlineMembers++
		}
	}

	GuildIconURL := "https://cdn.discordapp.com/icons/" + guild.ID + "/" + guild.Icon
	GuildCreationTime, err := utils.CreationTime(guild.ID)
	if err != nil {
		panic(err)
	}

	DaysAgo := strconv.Itoa(int(time.Since(GuildCreationTime).Hours() / 24))
	desc := "This server was created on " + GuildCreationTime.Format("Mon 01/02/2006") + " at " + GuildCreationTime.Format("03:04 pm") + ". That's over " + DaysAgo + " days ago!"

	var VerificationFmt string
	switch guild.VerificationLevel {
	case discord.VerificationLevelNone:
		VerificationFmt = "None"
	case discord.VerificationLevelLow:
		VerificationFmt = "Low"
	case discord.VerificationLevelMedium:
		VerificationFmt = "Medium"
	case discord.VerificationLevelHigh:
		VerificationFmt = "High"
	default:
		VerificationFmt = "Very High"
	}

	AnimatedEmojisCount := 0
	for _, e := range guild.Emojis {
		if e.Animated == true {
			AnimatedEmojisCount++
		}
	}

	TextChannels := 0
	VoiceChannels := 0

	for _, c := range guild.Channels {
		if c.Type == discord.ChannelTypeGuildText {
			TextChannels++
		} else if c.Type == discord.ChannelTypeGuildVoice {
			VoiceChannels++
		}
	}

	// gets the guild owner
	owner, err := ctx.Session.User(guild.OwnerID)
	if err != nil {
		return
	}

	em := utils.NewEmbed().
		SetAuthor(guild.Name).
		SetDescription(desc).
		SetColor(0x2ecc71).
		SetThumbnail(GuildIconURL).
		SetFooter("ID: " + guild.ID).
		AddField(utils.FieldParams{Name: "Owner", Value: owner.String()}).
		AddField(utils.FieldParams{Name: "Members", Value: strconv.Itoa(onlineMembers) + "/" + strconv.Itoa(guild.MemberCount) + " online"}).
		AddField(utils.FieldParams{Name: "Region", Value: guild.Region}).
		AddField(utils.FieldParams{Name: "Verification Level", Value: VerificationFmt}).
		AddField(utils.FieldParams{Name: "Text Channels", Value: strconv.Itoa(TextChannels)}).
		AddField(utils.FieldParams{Name: "Voice Channels", Value: strconv.Itoa(VoiceChannels)}).
		AddField(utils.FieldParams{Name: "Roles", Value: strconv.Itoa(len(guild.Roles))}).
		AddField(utils.FieldParams{Name: "Emojis", Value: strconv.Itoa(len(guild.Emojis)-AnimatedEmojisCount) + "/" + strconv.Itoa(AnimatedEmojisCount) + " (normal/animated)"}).
		InlineAllFields().MessageEmbed

	ctx.SendComplex("", em)
}

func userinfo(ctx *Context) {
	guild := ctx.Guild
	channel := ctx.Channel

	if channel.Type == discord.ChannelTypeDM || channel.Type == discord.ChannelTypeGroupDM {
		return
	}

	user := ctx.Message.Author

	AccountCreationTime, err := utils.CreationTime(user.ID)
	if err != nil {
		panic(err)
	}

	DaysAgo := strconv.Itoa(int(time.Since(AccountCreationTime).Hours() / 24))
	desc := "You created your account on " + AccountCreationTime.Format("Mon 01/02/2006") + " at " + AccountCreationTime.Format("03:04 pm") + ". That's over " + DaysAgo + " days ago!"

	members := guild.Members
	sort.Slice(members, func(i, j int) bool {
		a, err := members[i].JoinedAt.Parse()
		if err != nil {
			panic(err)
		}
		b, err := members[j].JoinedAt.Parse()
		if err != nil {
			panic(err)
		}
		return a.Unix() < b.Unix()
	})

	MemberCount := 0
	var member *discord.Member
	for i, v := range members {
		if v.User.ID == user.ID {
			member = v
			MemberCount = i + 1
		}
	}

	if member == nil {
		ctx.Send("User not found.")
		return
	}

	roles := utils.GetRoles(guild.Roles, member.Roles)
	RoleNames := make([]string, len(roles))
	for _, r := range roles {
		RoleNames = append(RoleNames, r.Name)
	}

	sort.Slice(roles, func(i, j int) bool {
		return roles[i].Position < roles[j].Position
	})

	RolesFmt := strings.Join(RoleNames, ", ")
	JoinedAt, err := member.JoinedAt.Parse()
	if err != nil {
		panic(err)
	}

	em := utils.NewEmbed().
		SetDescription(desc).
		SetColor(0x2ecc71).
		SetThumbnail(user.AvatarURL("")).
		SetAuthor([]string{user.String(), user.AvatarURL("")}...).
		SetFooter("ID: " + user.ID).
		AddField(utils.FieldParams{Name: "Name", Value: user.Username, Inline: true}).
		AddField(utils.FieldParams{Name: "Member No.", Value: strconv.Itoa(MemberCount), Inline: true}).
		AddField(utils.FieldParams{Name: "Account Created", Value: AccountCreationTime.Format("Mon 01/02/2006"), Inline: true}).
		AddField(utils.FieldParams{Name: "Joined At", Value: JoinedAt.Format("Mon 01/02/2006"), Inline: true}).
		AddField(utils.FieldParams{Name: "Roles", Value: RolesFmt, Inline: false}).
		MessageEmbed

	ctx.SendComplex("", em)
}

func init() {
	cog := NewCog("Info", "Information about certain things", false)
	cog.AddCommand("serverinfo", "Retrieves info about the server", []string{"si"}, serverinfo)
	cog.AddCommand("userinfo", "Gets info about a user", []string{"ui"}, userinfo)
	cog.Load()
}
