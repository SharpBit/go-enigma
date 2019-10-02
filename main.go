package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"./events"
	"./utils"

	discord "github.com/bwmarrin/discordgo"
)

func main() {
	bot, err := discord.New("Bot " + utils.GetConfig("token"))

	if err != nil {
		fmt.Println("error logging in,", err)
		return
	}

	bot.AddHandler(events.Ready)
	bot.AddHandler(events.MessageCreate)

	err = bot.Open()

	if err != nil {
		fmt.Println("error connecting to websocket,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	bot.Close()
}
