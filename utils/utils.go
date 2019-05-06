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
func GetRoles(roles []*discord.Role, ids []string) []*discord.Role {
	filter := make([]*discord.Role, len(ids))

	for _, r := range roles {
		for i, id := range ids {
			if id == r.ID {
				filter[i] = r
				break
			}
		}
	}

	return filter
}
