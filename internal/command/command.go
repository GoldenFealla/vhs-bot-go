package command

import (
	"github.com/bwmarrin/discordgo"
)

type Slash interface {
	Data() *discordgo.ApplicationCommand
	Handler(s *discordgo.Session, i *discordgo.InteractionCreate) error
}

var Slashes = map[string]Slash{
	"server": ServerSlashCommand(),
	"user":   UserSlashCommand(),
}
