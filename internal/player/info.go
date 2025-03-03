package player

import (
	"bytes"
	"encoding/json"
	"errors"
	URL "net/url"
	"os/exec"
)

type Extractor = string

var (
	YOUTUBE  Extractor = "youtube"
	BANDCAMP Extractor = "bandcamp"
)

type VideoData struct {
	ID             string
	URL            string
	Title          string
	Channel        string
	ChannelURL     string
	Thumbnail      string
	DurationString string
	StreamURL      string
}

func MusicInfo(url string) (*VideoData, error) {
	_, err := URL.ParseRequestURI(url)
	if err != nil {
		return nil, err
	}

	args := []string{
		"--extractor-args",
		"youtube:skip=hls,dash,translated_subs",
		"--ignore-no-formats-error",
		"--skip-download",
		"--dump-single-json",
		"--flat-playlist",
		"--parse-metadata", "\":(?P<thumbnails>)\"",
		"--parse-metadata", "\":(?P<automatic_captions>)\"",
		"--parse-metadata", "\":(?P<heatmap>)\"",
		url,
	}

	// dw: stdout, ew: stderr
	var dw bytes.Buffer
	var ew bytes.Buffer

	command := exec.Command("yt-dlp", args...)
	command.Stderr = &ew
	command.Stdout = &dw

	err = command.Start()
	if err != nil {
		return nil, err
	}

	err = command.Wait()
	if err != nil {
		return nil, err
	}

	var data map[string]any
	if err := json.Unmarshal(dw.Bytes(), &data); err != nil {
		return nil, err
	}

	if ew.Len() > 0 {
		return nil, errors.New(ew.String())
	}

	var streamURL string

	switch data["extractor"].(string) {
	case YOUTUBE:
		streamURL = data["requested_formats"].([]any)[1].(map[string]any)["url"].(string)
	case BANDCAMP:
		streamURL = data["requested_downloads"].([]any)[1].(map[string]any)["url"].(string)
	}

	VideoData := &VideoData{
		ID:             data["id"].(string),
		URL:            data["webpage_url"].(string),
		Title:          data["title"].(string),
		Channel:        data["channel"].(string),
		ChannelURL:     data["channel_url"].(string),
		Thumbnail:      data["thumbnail"].(string),
		DurationString: data["duration_string"].(string),
		StreamURL:      streamURL,
	}

	return VideoData, nil
}
