package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	"goldenfealla/vhs-bot/config"
	"goldenfealla/vhs-bot/handler"
	"goldenfealla/vhs-bot/internal/command"
)

var (
	RemoveCommand   = flag.Bool("remove", false, "Remove All Command")
	RegisterCommand = flag.Bool("register", false, "Register All Command")
)

func registerCommands(s *discordgo.Session) {
	log.Println("Adding commands...")
	for _, v := range command.Slashes {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", v.Data)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Data.Name, err)
		}
	}
}

func removeCommands(s *discordgo.Session) {
	log.Println("Removing commands...")

	registeredCommands, err := s.ApplicationCommands(s.State.User.ID, "")
	if err != nil {
		log.Fatalf("Could not fetch registered commands: %v", err)
	}

	log.Println("Removing commands...")
	for _, v := range registeredCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, "", v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
}

func run() {
	// OS Interupt
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	<-sigch
}

func init() {
	flag.Parse()
}

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
	session.AddHandler(handler.InteractionCreate)

	err = session.Open()
	if err != nil {
		log.Fatalf("could not open session: %s", err)
	}

	if *RegisterCommand {
		registerCommands(session)
	} else if *RemoveCommand {
		removeCommands(session)
	} else {
		run()
	}

	err = session.Close()
	if err != nil {
		log.Printf("could not close session gracefully: %s", err)
	}
}
