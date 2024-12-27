package command

import (
	"goldenfealla/vhs-bot/config"

	"github.com/bwmarrin/discordgo"
)

const TIME_FORMAT = " 15:04:05 02/01/2006"

func UserSlashCommand() *SlashCommand {
	data := &discordgo.ApplicationCommand{
		Name:        "user",
		Description: "user info",
	}

	handler := func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		user := i.Member

		// user id always available since it come from caller
		createdAt, _ := discordgo.SnowflakeTimestamp(user.User.ID)

		embed := discordgo.MessageEmbed{
			Title: "User basic information",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  user.DisplayName(),
					Value: user.User.Username,
				},
				{
					Name:  "Created At",
					Value: createdAt.Local().Format(TIME_FORMAT),
				},
				{
					Name:  "Joined this server At",
					Value: user.JoinedAt.Local().Format(TIME_FORMAT),
				},
			},
			Image: &discordgo.MessageEmbedImage{
				URL: user.AvatarURL("1024"),
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
