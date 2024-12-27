package command

import (
	"github.com/bwmarrin/discordgo"
)

type SlashCommand struct {
	Data    *discordgo.ApplicationCommand
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var Slashes = map[string]*SlashCommand{
	"server": ServerSlashCommand(),
	"user":   UserSlashCommand(),
}
