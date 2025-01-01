package command

import (
	"fmt"

	"goldenfealla/vhs-bot/config"

	"github.com/bwmarrin/discordgo"
)

type serverSlash struct{}

func ServerSlashCommand() *serverSlash {
	return &serverSlash{}
}

func (sc serverSlash) Data() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "server",
		Description: "Server info",
	}
}

func (sc serverSlash) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	guild, err := s.GuildWithCounts(i.GuildID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Error: %v", err.Error()),
			},
		})
		return err
	}

	totalMembers := guild.ApproximateMemberCount
	presenceMembers := guild.ApproximatePresenceCount

	embed := discordgo.MessageEmbed{
		Title: "Server asic information",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: fmt.Sprintf("ID: %v", guild.ID),
			},
			{
				Name:  guild.Name,
				Value: fmt.Sprintf("Owner: <@!%v>", guild.OwnerID),
			},
			{
				Name:  fmt.Sprintf("Members: %v", totalMembers),
				Value: fmt.Sprintf("Online %v | Offline: %v", presenceMembers, totalMembers-presenceMembers),
			},
		},
		Image: &discordgo.MessageEmbedImage{
			URL: guild.IconURL("256"),
		},
		Color: config.DEFAULT_COLOR,
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{&embed},
		},
	})
}
