package utils

import (
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
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
func GetRoles(guild *discordgo.Guild, member *discordgo.Member) []*discordgo.Role {
	roles := []*discordgo.Role{}
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

// CleanupCode removes markdown around the code
func CleanupCode(content string) (code string) {
	if strings.HasPrefix(content, "```") && strings.HasSuffix(content, "```") {
		lines := strings.Split(content, "\n")
		return strings.Join(lines[1:len(lines)-1], "\n")
	}
	return strings.Trim(content, "` \n")
}
