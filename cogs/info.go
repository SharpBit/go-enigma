package cogs

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/SharpBit/go-enigma/commands"
	"github.com/SharpBit/go-enigma/utils"

	discord "github.com/bwmarrin/discordgo"
)

func serverinfo(ctx *commands.Context) (err error) {
	guild := ctx.Guild
	channel := ctx.Channel

	if channel.Type == discord.ChannelTypeDM || channel.Type == discord.ChannelTypeGroupDM {
		return fmt.Errorf("CommandCheckError: Command cannot be run in a DM.")
	}

	// Get the number of online members
	onlineMembers := 0
	for _, m := range guild.Presences {
		if m.Status != discord.StatusOffline {
			onlineMembers++
		}
	}

	// Build the Icon URL and get the time of guild creation
	GuildIconURL := "https://cdn.discordapp.com/icons/" + guild.ID + "/" + guild.Icon
	GuildCreationTime, err := utils.CreationTime(guild.ID)
	if err != nil {
		return
	}

	// Format the description
	DaysAgo := strconv.Itoa(int(time.Since(GuildCreationTime).Hours() / 24))
	desc := "This server was created on " + GuildCreationTime.Format("Mon 01/02/2006") + " at " + GuildCreationTime.Format("03:04 pm") + ". That's over " + DaysAgo + " days ago!"

	// Convert constants into human readable verification levels
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

	// Get the number of Animated emojis
	AnimatedEmojisCount := 0
	for _, e := range guild.Emojis {
		if e.Animated == true {
			AnimatedEmojisCount++
		}
	}

	// Get the number of text and voice channels
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

	// Build the embed
	em := commands.NewEmbed().
		SetAuthor(guild.Name).
		SetDescription(desc).
		SetColor(0x2ecc71).
		SetThumbnail(GuildIconURL).
		SetFooter("ID: "+guild.ID).
		AddField("Owner", owner.String()).
		AddField("Members", strconv.Itoa(onlineMembers)+"/"+strconv.Itoa(guild.MemberCount)+" online").
		AddField("Region", guild.Region).
		AddField("Verification Level", VerificationFmt).
		AddField("Text Channels", strconv.Itoa(TextChannels)).
		AddField("Voice Channels", strconv.Itoa(VoiceChannels)).
		AddField("Roles", strconv.Itoa(len(guild.Roles))).
		AddField("Emojis", strconv.Itoa(len(guild.Emojis)-AnimatedEmojisCount)+"/"+strconv.Itoa(AnimatedEmojisCount)+" (normal/animated)").
		InlineAllFields().MessageEmbed

	_, err = ctx.SendComplex("", em)
	return
}

func userinfo(ctx *commands.Context, member *discord.Member) (err error) {
	guild := ctx.Guild
	channel := ctx.Channel

	if channel.Type == discord.ChannelTypeDM || channel.Type == discord.ChannelTypeGroupDM {
		return fmt.Errorf("CommandCheckError: Command cannot be run in a DM.")
	}

	var user *discord.User
	if member == nil {
		member = ctx.Message.Member
		user = ctx.Author
	} else {
		user = member.User
	}

	AccountCreationTime, err := utils.CreationTime(user.ID)
	if err != nil {
		return
	}

	DaysAgo := strconv.Itoa(int(time.Since(AccountCreationTime).Hours() / 24))
	desc := "You created your account on " + AccountCreationTime.Format("Mon 01/02/2006") + " at " + AccountCreationTime.Format("03:04 pm") + ". That's over " + DaysAgo + " days ago!"

	members := guild.Members
	sort.Slice(members, func(i, j int) bool {
		a, _ := members[i].JoinedAt.Parse()
		b, _ := members[j].JoinedAt.Parse()
		return a.Unix() < b.Unix()
	})

	// Get the Member Join Position
	var MemberCount int
	for i, v := range members {
		if v.User.ID == user.ID {
			MemberCount = i + 1
		}
	}

	// Sort the roles from Highest to Lowest
	roles := utils.GetRoles(guild, member)
	sort.Slice(roles, func(i, j int) bool {
		return roles[i].Position > roles[j].Position
	})

	// Get the role names from the role list
	RoleNames := []string{}
	for _, r := range roles {
		RoleNames = append(RoleNames, r.Name)
	}

	// Format the list of roles and get the join time
	RolesFmt := strings.Join(RoleNames, ", ")
	JoinedAt, err := member.JoinedAt.Parse()
	if err != nil {
		return
	}

	// Build the embed
	em := commands.NewEmbed().
		SetDescription(desc).
		SetColor(0x2ecc71).
		SetThumbnail(user.AvatarURL("256")).
		SetAuthor(user.String(), user.AvatarURL("64")).
		SetFooter("ID: "+user.ID).
		AddField("Name", user.Username, true).
		AddField("Member No.", strconv.Itoa(MemberCount), true).
		AddField("Account Created", AccountCreationTime.Format("Mon 01/02/2006"), true).
		AddField("Joined At", JoinedAt.Format("Mon 01/02/2006"), true).
		AddField("Roles", RolesFmt, false).
		MessageEmbed

	_, err = ctx.SendComplex("", em)
	return
}

func avatar(ctx *commands.Context, user *discord.User) (err error) {
	if user == nil {
		user = ctx.Message.Author
	}

	em := commands.NewEmbed().
		SetColor(0x2ecc71).
		SetAuthor(user.String(), user.AvatarURL("64")).
		SetFooter("ID: " + user.ID).
		SetImage(user.AvatarURL("2048")).
		MessageEmbed

	_, err = ctx.SendComplex("", em)
	return
}

func servericon(ctx *commands.Context) (err error) {
	if ctx.Channel.Type == discord.ChannelTypeDM || ctx.Channel.Type == discord.ChannelTypeGroupDM {
		return fmt.Errorf("CommandCheckError: Command cannot be run in a DM.")
	}

	guild := ctx.Guild

	em := commands.NewEmbed().
		SetColor(0x2ecc71).
		SetAuthor(guild.Name, guild.IconURL()+"?size=64").
		SetFooter("ID: " + guild.ID).
		SetImage(guild.IconURL() + "?size=2048").
		MessageEmbed

	_, err = ctx.SendComplex("", em)
	return err
}

func init() {
	cog := commands.NewCog("Info", "Information about certain things")
	cog.AddCommand("serverinfo", "Retrieves info about the server", "", serverinfo).
		SetAliases("server", "si").
		AddCheck(utils.GuildOnly)
	cog.AddCommand("userinfo", "Gets info about a user", "[member]", userinfo).
		SetAliases("user", "ui").
		SetDefaultArg(nil).
		AddCheck(utils.GuildOnly)
	cog.AddCommand("avatar", "Get the avatar for a certain user", "[user]", avatar).
		SetAliases("av").
		SetDefaultArg(nil)
	cog.AddCommand("servericon", "Get the server icon", "", servericon).
		SetAliases("icon").
		AddCheck(utils.GuildOnly)
	cog.Load()
}
