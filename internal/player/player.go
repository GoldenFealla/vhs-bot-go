package player

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/kkdai/youtube/v2"
)

var client = youtube.Client{}

func Play(s *discordgo.Session, i *discordgo.InteractionCreate, url string) error {
	// Join the provided voice channel.
	g, err := s.State.Guild(i.GuildID)
	if err != nil {
		// Could not find guild.
		return fmt.Errorf("no guild found")
	}

	cID := ""

	for _, vs := range g.VoiceStates {
		if vs.UserID == i.Member.User.ID {
			cID = vs.ChannelID
			break
		}
		return fmt.Errorf("user not in voice")
	}

	vc, err := s.ChannelVoiceJoin(g.ID, cID, false, true)
	if err != nil {
		return fmt.Errorf("Error joining: %v", err)
	}

	// Sleep for a specified amount of time before playing the sound
	time.Sleep(250 * time.Millisecond)

	// process video
	video, err := client.GetVideo(url)
	if err != nil {
		return err
	}

	audioFormats := video.Formats.Type("audio/mp4")

	_, err = client.GetStreamURL(video, &audioFormats[len(audioFormats)-1])
	if err != nil {
		return err
	}

	// Start speaking.
	vc.Speaking(true)

	vc.Speaking(false)

	// Sleep for a specificed amount of time before ending.
	time.Sleep(250 * time.Millisecond)
	return nil
}
