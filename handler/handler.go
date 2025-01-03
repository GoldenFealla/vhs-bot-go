package handler

import (
	"goldenfealla/vhs-bot/internal/command"
	"log"

	"github.com/bwmarrin/discordgo"
)

func Ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Printf("Logged in as: %v#%v", event.User.Username, event.User.Discriminator)
}

func InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	name := i.ApplicationCommandData().Name

	if slash, ok := command.Slashes[name]; ok {
		err := slash.Handler(s, i)
		if err != nil {
			log.Printf("cmd: %v, err: %v", name, err)
		}
	}
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	// No interaction now
}
