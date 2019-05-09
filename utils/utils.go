package utils

import (
	"strconv"
	"time"

	discord "github.com/bwmarrin/discordgo"
)

// CreationTime returns the creation time of a Snowflake ID relative to the creation of Discord.
func CreationTime(ID string) (t time.Time, err error) {
	i, err := strconv.ParseInt(ID, 10, 64)
	if err != nil {
		return
	}
	timestamp := (i >> 22) + 1420070400000
	t = time.Unix(timestamp/1000, 0)
	return
}

// GetRoles finds the roles based on given roles and IDs
func GetRoles(guild *discord.Guild, member *discord.Member) []*discord.Role {
	roles := []*discord.Role{}
	for _, m := range member.Roles {
		for _, r := range guild.Roles {
			if r.ID == m {
				roles = append(roles, r)
				break
			}
		}
	}
	return roles
}
