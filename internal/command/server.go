package command

import (
	"fmt"
	"log"

	"goldenfealla/vhs-bot/config"

	"github.com/bwmarrin/discordgo"
)

func ServerSlashCommand() *SlashCommand {
	data := &discordgo.ApplicationCommand{
		Name:        "server",
		Description: "Server info",
	}

	handler := func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		guild, err := s.GuildWithCounts(i.GuildID)
		if err != nil {
			log.Println(err)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Error: %v", err.Error()),
				},
			})
			return
		}

		totalMembers := guild.ApproximateMemberCount
		presenceMembers := guild.ApproximatePresenceCount

		embed := discordgo.MessageEmbed{
			Title: "Server basic information",
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

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{&embed},
			},
		})
	}

	return &SlashCommand{
		Data:    data,
		Handler: handler,
	}
}
