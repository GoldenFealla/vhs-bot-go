package command

import (
	"fmt"
	"goldenfealla/vhs-bot/config"
	"goldenfealla/vhs-bot/internal/player"

	"github.com/bwmarrin/discordgo"
)

type playSlash struct{}

func PlaySlashCommand() *playSlash {
	return &playSlash{}
}

func (sc playSlash) Data() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "play",
		Description: "Extract play info",
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

func (sc playSlash) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	url := optionMap["url"].StringValue()
	DeferReply(s, i)

	if player.Players[i.GuildID] == nil {
		player.Players[i.GuildID] = player.NewPlayer()
	}

	g, err := s.State.Guild(i.GuildID)
	if err != nil {
		EditReplyString(s, i, err.Error())
		return err
	}

	channelID := ""
	isFound := false

	for _, vs := range g.VoiceStates {
		if vs.UserID == i.Member.User.ID {
			channelID = vs.ChannelID
			isFound = true
			break
		}
	}

	if !isFound {
		EditReplyString(s, i, "no channel found or user is not in voice")
		return fmt.Errorf("no channel found or user is not in voice")
	}

	vi, err := player.Players[i.GuildID].Play(s, g.ID, channelID, url)
	if err != nil {
		EditReplyString(s, i, err.Error())
		return err
	}

	embed := &discordgo.MessageEmbed{
		URL:         vi.URL,
		Title:       vi.Title,
		Description: fmt.Sprintf("**Channel: [%v](%v)**", vi.Channel, vi.ChannelURL),
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
