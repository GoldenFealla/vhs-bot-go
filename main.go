package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	"goldenfealla/vhs-bot/config"
	"goldenfealla/vhs-bot/handler"
)

func main() {
	// config
	config.Init()

	// add env here
	godotenv.Load()

	token := os.Getenv("DISCORD_TOKEN")

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Println(err)
	}

	session.AddHandler(handler.Ready)
	session.AddHandler(handler.MessageCreate)

	err = session.Open()
	if err != nil {
		log.Fatalf("could not open session: %s", err)
	}

	// OS Interupt
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	<-sigch

	err = session.Close()
	if err != nil {
		log.Printf("could not close session gracefully: %s", err)
	}
}
