package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type ChannelItems struct {
	Items []struct {
		ContentDetails struct {
			RelatedPlaylists struct {
				Uploads string `json:"uploads"`
			} `json:"relatedPlaylists"`
		} `json:"contentDetails"`
	} `json:"items"`
}

type PlaylistItems struct {
	Items []struct {
		Snippet struct {
			Title      string `json:"title"`
			ResourceId struct {
				VideoId string `json:"videoId"`
			} `json:"resourceId"`
		} `json:"snippet"`
	} `json:"items"`
}

// clean the title
func cleanTitle(title string) string {
	// Remove hashtags and the text immediately following them
	re := regexp.MustCompile(`#\S+`)
	cleaned := re.ReplaceAllString(title, "")
	// Remove extra spaces
	cleaned = strings.TrimSpace(cleaned)
	cleaned = regexp.MustCompile(`\s+`).ReplaceAllString(cleaned, " ")
	return cleaned
}

// Get the playlist id
func getPlaylistId(apiKey, channelHandle string) (string, error) {
	url := fmt.Sprintf("https://youtube.googleapis.com/youtube/v3/channels?part=contentDetails&forHandle=" + channelHandle + "&key=" + apiKey + "&maxResults=5")
	channelRes, err := http.Get(url)

	// if invalid request
	if err != nil {
		return "", err
	}
	defer channelRes.Body.Close()
	if channelRes.StatusCode != 200 {
		fmt.Println("Youtube Channels API not Available")
	}

	// get response body
	channelsBody, err := io.ReadAll(channelRes.Body)
	if err != nil {
		return "", err
	}

	// unmarshal channels
	var channel ChannelItems
	if err := json.Unmarshal(channelsBody, &channel); err != nil {
		return "", err
	}
	return channel.Items[0].ContentDetails.RelatedPlaylists.Uploads, nil
}

func getPlaylistItems(apiKey, playlistId string) (title, videoId string, err error) {
	url := fmt.Sprintf("https://youtube.googleapis.com/youtube/v3/playlistItems?part=snippet&key=%s&playlistId=%s&maxResults=1", apiKey, playlistId)
	playlistRes, err := http.Get(url)
	if err != nil {
		return "", "", err
	}

	defer playlistRes.Body.Close()
	if playlistRes.StatusCode != 200 {
		fmt.Println("Youtube PlaylistItems API not Available")
	}

	// get body
	playlistItemsBody, err := io.ReadAll(playlistRes.Body)
	if err != nil {
		return "", "", err
	}

	// ummarshal
	var playlistItems PlaylistItems

	json.Unmarshal(playlistItemsBody, &playlistItems)
	title = playlistItems.Items[0].Snippet.Title
	videoId = playlistItems.Items[0].Snippet.ResourceId.VideoId
	title = cleanTitle(title)
	return title, videoId, nil

}

func main() {
	channelHandle := "triplescosmos"
	// load env, make sure you export api key
	apiKey := os.Getenv("YOUTUBE_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "ERROR: YOUTUBE_API_KEY is not set")
		os.Exit(1)
	}

	uploadsPlaylistId, err := getPlaylistId(apiKey, channelHandle)
	if err != nil {
		panic(err)
	}

	videoTitle, videoId, err := getPlaylistItems(apiKey, uploadsPlaylistId)
	if err != nil {
		panic(err)
	}

	fmt.Println(videoTitle)
	fmt.Printf("Enjoy Your Content! https://www.youtube.com/watch?v=%s\n", videoId)
}
