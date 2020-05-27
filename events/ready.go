package events

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Ready event to update status on startup
func Ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateStatus(0, "Go, Discord, go!")
	fmt.Println("Bot is ready to GO!")
}
