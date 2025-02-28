package player

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"
)

type Format struct {
	Url string `json:"url"`
}

type VideoData struct {
	ID               string   `json:"id"`
	Title            string   `json:"title"`
	Thumbnail        string   `json:"thumbnail"`
	DurationString   string   `json:"duration_string"`
	RequestedFormats []Format `json:"requested_formats"`
}

func fetch(url string, w io.WriteCloser) error {
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
	command.Stderr = os.Stderr
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
	var data VideoData

	r, w := io.Pipe()
	defer r.Close()

	go fetch(url, w)

	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}
