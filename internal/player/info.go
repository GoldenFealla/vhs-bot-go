package player

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	URL "net/url"
	"os/exec"
)

type Extractor = string

var (
	YOUTUBE  Extractor = "Youtube"
	BANDCAMP Extractor = "Bandcamp"
)

type Format struct {
	Url string `json:"url"`
}

type Download struct {
	RequestedFormats []Format `json:"requested_formats"`
	Url              string   `json:"url"`
}

type VideoData struct {
	ID                 string     `json:"id"`
	Title              string     `json:"title"`
	Thumbnail          string     `json:"thumbnail"`
	DurationString     string     `json:"duration_string"`
	Extractor          Extractor  `json:"extractor"`
	RequestedDownloads []Download `json:"requested_downloads"`
}

func fetch(url string, w io.WriteCloser, errWriter io.Writer) error {
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

	command := exec.Command("yt-dlp", args...)
	command.Stderr = errWriter
	command.Stdout = w
	defer w.Close()

	err := command.Start()
	if err != nil {
		return err
	}

	err = command.Wait()
	if err != nil {
		return err
	}

	return nil
}

func Info(url string) (*VideoData, error) {
	_, err := URL.ParseRequestURI(url)
	if err != nil {
		return nil, err
	}

	var data VideoData

	r, w := io.Pipe()
	defer r.Close()

	var errBuf bytes.Buffer

	go func() {
		err := fetch(url, w, &errBuf)
		if err != nil {
			log.Println(err)
		}
	}()

	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}

	if errBuf.Len() > 0 {
		return nil, errors.New(errBuf.String())
	}

	return &data, nil
}
