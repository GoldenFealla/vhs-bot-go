package command

import (
	"goldenfealla/vhs-bot/internal/player"

	"github.com/bwmarrin/discordgo"
)

type youtubeSlash struct{}

func YoutubeSlashCommand() *youtubeSlash {
	return &youtubeSlash{}
}

func (sc youtubeSlash) Data() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "youtube",
		Description: "Extract youtube info",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "url",
				Description: "This is the url",
				Required:    true,
			},
		},
	}
}

func (sc youtubeSlash) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	url := optionMap["url"].StringValue()

	go player.Play(s, i, url)

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Done",
		},
	})
}
