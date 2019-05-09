package commands

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"../utils"

	discord "github.com/bwmarrin/discordgo"
)

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
	cmd, existing := NewCommand("userinfo", "Retrieves info about a user")
	if existing {
		fmt.Println("error: command userinfo already exists")
		return
	}
	cmd.Run = userinfo
	cmd.Aliases = []string{"ui"}
	RegisterCommand(cmd)
}
