package player

import (
	"testing"
)

func TestMusicInfoYoutube(t *testing.T) {
	musicInfo, err := MusicInfo("https://www.youtube.com/watch?v=q6i72LB9kT4")
	if err != nil {
		t.Fatal("Error get music info")
	}

	wantedID := "q6i72LB9kT4"
	wantedURL := "https://www.youtube.com/watch?v=q6i72LB9kT4"
	wantedTitle := "勾指起誓【翻唱 ▪ 泠鳶yousa】ilem&洛天依"
	wantedChannel := "泠鳶yousa【Unofficial Channel】"
	wantedChannelURL := "https://www.youtube.com/channel/UCT64yx_dtZvlRCU3xYU-alg"

	t.Run("Music Info Match ID", func(t *testing.T) {
		if musicInfo.ID != wantedID {
			t.Errorf("ID value: %q, Expected: %#q", musicInfo.ID, wantedID)
		}
	})

	t.Run("Music Info Match URL", func(t *testing.T) {
		if musicInfo.URL != wantedURL {
			t.Errorf("URL value: %q, Expected: %#q", musicInfo.URL, wantedURL)
		}
	})

	t.Run("Music Info Match Title", func(t *testing.T) {
		if musicInfo.Title != wantedTitle {
			t.Errorf("Title value: %q, Expected: %#q", musicInfo.Title, wantedTitle)
		}
	})

	t.Run("Music Info Match Channel", func(t *testing.T) {
		if musicInfo.Channel != wantedChannel {
			t.Errorf("Channel value: %q, Expected: %#q", musicInfo.Channel, wantedChannel)
		}
	})

	t.Run("Music Info Match Channel URL", func(t *testing.T) {
		if musicInfo.ChannelURL != wantedChannelURL {
			t.Errorf("ChannelURL value: %q, Expected: %#q", musicInfo.ChannelURL, wantedChannelURL)
		}
	})

	t.Run("Music Info Stream URL not empty", func(t *testing.T) {
		if musicInfo.StreamURL == "" {
			t.Errorf("StreamURL value: %q, Expected: not empty", musicInfo.ChannelURL)
		}
	})
}
