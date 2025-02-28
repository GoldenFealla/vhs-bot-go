package command

import (
	"fmt"
	"goldenfealla/vhs-bot/config"
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
				Description: "This is the url needed to extract",
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

	DeferReply(s, i)

	g, err := s.State.Guild(i.GuildID)
	if err != nil {
		EditReplyString(s, i, err.Error())
		return err
	}

	channelID, err := player.GetChannelID(g, i.Member.User.ID)
	if err != nil {
		EditReplyString(s, i, err.Error())
		return err
	}

	vc, err := player.Join(s, g.ID, channelID)
	if err != nil {
		EditReplyString(s, i, err.Error())
		return err
	}

	vi, err := player.Play(vc, url)
	if err != nil {
		EditReplyString(s, i, err.Error())
		return err
	}

	embed := &discordgo.MessageEmbed{
		Title: vi.Title,
		Image: &discordgo.MessageEmbedImage{
			URL: vi.Thumbnail,
		},
		Color: config.DEFAULT_COLOR,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Duration: %v", vi.DurationString),
		},
	}

	EditReplyEmbed(s, i, embed)
	return nil
}
