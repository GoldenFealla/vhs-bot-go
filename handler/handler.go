package handler

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func Ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Printf("Logged in as: %v#%v", event.User.Username, event.User.Discriminator)
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Message received: %s", m.Content))
}
