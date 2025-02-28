package player

import (
	"errors"
	"fmt"
	"goldenfealla/vhs-bot/internal/encode"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func process(vc *discordgo.VoiceConnection, vi *VideoData) {
	vc.Speaking(true)

	outputChan := make(chan []byte, 64)

	go func() {
		var err error

		if len(vi.RequestedDownloads) == 0 {
			err = errors.New("something wrong")
		} else {
			switch vi.Extractor {
			case YOUTUBE:
				err = encode.Encode(vi.RequestedDownloads[0].RequestedFormats[1].Url, outputChan)
			case BANDCAMP:
				err = encode.Encode(vi.RequestedDownloads[0].Url, outputChan)
			default:
				err = errors.New("not supported site")
			}
		}

		if err != nil {
			log.Println(err)
		}
	}()

	for data := range outputChan {
		vc.OpusSend <- data
	}

	vc.Speaking(false)
	time.Sleep(250 * time.Millisecond)

	vc.Disconnect()
}

func GetChannelID(g *discordgo.Guild, userID string) (channelID string, err error) {
	isFound := false

	for _, vs := range g.VoiceStates {
		if vs.UserID == userID {
			channelID = vs.ChannelID
			isFound = true
			break
		}
	}

	if isFound {
		return channelID, nil
	}

	return "", fmt.Errorf("no channel found")
}

func Join(s *discordgo.Session, guildID string, channelID string) (*discordgo.VoiceConnection, error) {
	return s.ChannelVoiceJoin(guildID, channelID, false, true)
}

func Play(vc *discordgo.VoiceConnection, url string) (*VideoData, error) {
	vi, err := Info(url)
	if err != nil {
		return nil, fmt.Errorf("error get video info: %v", err)
	}

	go process(vc, vi)

	return vi, nil
}
