package events

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"../utils"

	discord "github.com/bwmarrin/discordgo"
)

var (
	prefix = "//"
	owner  = "281821029490229251"
)

// MessageCreate event
func MessageCreate(session *discord.Session, msg *discord.MessageCreate) {

	// so the bot doesn't respond to other bots or webhooks (including itself)
	if msg.Author.Bot || msg.WebhookID != "" {
		return
	}

	if !strings.HasPrefix(msg.Content, prefix) {
		return
	}

	// Separates the commands from the arguments
	input := strings.Split(msg.Content, " ")
	_, args := strings.Trim(input[0], prefix), input[1:]

	// cmd, ok := utils.CommandMap[CmdString]
	// if !ok {
	// 	return
	// }

	// gets the message's guild
	guild, err := session.State.Guild(msg.GuildID)
	if err != nil {
		guild, err = session.Guild(msg.GuildID)
		if err != nil {
			return
		}
	}

	GuildIconURL := "https://cdn.discordapp.com/icons/" + guild.ID + "/" + guild.Icon

	// gets the message's channel
	channel, err := session.State.Channel(msg.ChannelID)
	if err != nil {
		channel, err = session.Channel(msg.ChannelID)
		if err != nil {
			return
		}
	}

	// cmd(&commands.Context{
	// 	Session: session,
	// 	Message: msg,
	// 	Guild:   guild,
	// 	Channel: channel,
	// 	Author:  msg.Author,
	// 	Args:    args,
	// 	Prefix:  prefix,
	// })
	// return

	// ping command
	if msg.Content == prefix+"ping" {
		session.ChannelMessageSend(msg.ChannelID, "Pong!")
	}

	/*
	   INFO COMMANDS
	*/

	// serverinfo command
	if msg.Content == prefix+"serverinfo" || msg.Content == prefix+"si" {
		if channel.Type == discord.ChannelTypeDM || channel.Type == discord.ChannelTypeGroupDM {
			return
		}

		onlineMembers := 0
		for _, m := range guild.Presences {
			if m.Status != discord.StatusOffline {
				onlineMembers++
			}
		}

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
		owner, err := session.User(guild.OwnerID)
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

		session.ChannelMessageSendEmbed(msg.ChannelID, em)
	}

	// // userinfo command
	if msg.Content == prefix+"userinfo" || msg.Content == prefix+"ui" {
		if channel.Type == discord.ChannelTypeDM || channel.Type == discord.ChannelTypeGroupDM {
			return
		}

		user := msg.Author

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
			session.ChannelMessageSend(channel.ID, "User not found.")
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

		session.ChannelMessageSendEmbed(msg.ChannelID, em)
	}

	/*
	   OWNER COMMANDS
	*/

	// presence command
	if strings.HasPrefix(msg.Content, prefix+"presence") && msg.Author.ID == "281821029490229251" {
		switch args[0] {
		case "play":
			session.UpdateStatus(0, strings.Join(args[1:], " "))
			session.ChannelMessageSend(msg.ChannelID, "Set `play` status to `"+strings.Join(args[1:], " ")+"`")
		case "listen":
			session.UpdateListeningStatus(strings.Join(args[1:], " "))
			session.ChannelMessageSend(msg.ChannelID, "Set `listen` status to `"+strings.Join(args[1:], " ")+"`")
		case "stream":
			CleanURL := strings.Replace(strings.Replace(args[1], "<", "", 1), ">", "", 1)
			session.UpdateStreamingStatus(0, strings.Join(args[2:], " "), CleanURL)
			session.ChannelMessageSend(msg.ChannelID, "Set `stream` status to `"+strings.Join(args[2:], " ")+"`"+" at URL <"+CleanURL+">")
		case "reset":
			session.UpdateStatus(0, "")
		default:
			session.ChannelMessageSend(msg.ChannelID, args[0]+" is not an option. Select from `play`, `listen`, `stream`, `reset`")
		}
	}
}
