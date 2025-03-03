package player

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Player struct {
	c          *Controller
	vc         *discordgo.VoiceConnection
	tracks     []*VideoData
	cur        int
	isLoop     bool
	isFinished bool
	next       chan bool
}

func (p *Player) listenNext() {
	for range p.next {
		p.c.StopStream()
		if p.cur >= len(p.tracks)-1 {
			p.c.isPlaying = false
			p.isFinished = true
			p.vc.Speaking(false)
			continue
		}

		p.cur += 1
		p.c.Stream(p.tracks[p.cur].StreamURL, p.vc.OpusSend, p.next)
	}
}

func NewPlayer() *Player {
	player := &Player{
		c:          &Controller{isPlaying: false},
		vc:         nil,
		tracks:     []*VideoData{},
		cur:        -1,
		isLoop:     false,
		isFinished: true,
		next:       make(chan bool),
	}

	go player.listenNext()
	return player
}

var Players = map[string]*Player{}

// Controller
func (p *Player) Join(s *discordgo.Session, guildID string, channelID string) error {
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}
	p.vc = vc
	return nil
}

func (p *Player) Play(s *discordgo.Session, guildID string, channelID string, url string) (*VideoData, error) {
	vd, err := p.Queue(url)
	if err != nil {
		return nil, err
	}

	if p.vc == nil {
		p.Join(s, guildID, channelID)
	}

	if !p.c.isPlaying && p.isFinished {
		p.cur += 1
		p.vc.Speaking(true)
		p.isFinished = false
		go p.c.Stream(vd.StreamURL, p.vc.OpusSend, p.next)
	}

	return vd, nil
}

func (p *Player) Pause() {
	p.c.isPlaying = false
}

func (p *Player) Resume() {
	p.c.isPlaying = true
}

func (p *Player) Skip() error {
	p.c.StopStream()

	if p.cur >= len(p.tracks)-1 {
		return errors.New("no track left")
	}

	p.cur += 1
	p.c.Stream(p.tracks[p.cur].StreamURL, p.vc.OpusSend, p.next)
	return nil
}

func (p *Player) SkipTo(pos int) error {
	p.c.StopStream()

	if pos+1 >= len(p.tracks)-1 {
		return errors.New("invalid position")
	}

	p.cur = pos + 1
	p.c.Stream(p.tracks[p.cur].StreamURL, p.vc.OpusSend, p.next)
	return nil
}

func (p *Player) Loop() bool {
	p.isLoop = !p.isLoop
	return p.isLoop
}

func (p *Player) Leave() error {
	err := p.vc.Disconnect()

	p.vc = nil
	p.tracks = []*VideoData{}
	p.cur = 0

	return err
}

// Info
func (p *Player) Queue(url string) (*VideoData, error) {
	vd, err := MusicInfo(url)
	if err != nil {
		return nil, fmt.Errorf("error get video info: %v", err)
	}

	p.tracks = append(p.tracks, vd)
	return vd, nil
}

func (p *Player) List() []*VideoData {
	return p.tracks
}
