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

func cleanTitle(title string) string {
	// Remove hashtags and the text immediately following them
	re := regexp.MustCompile(`#\S+`)
	cleaned := re.ReplaceAllString(title, "")

	// Remove extra spaces
	cleaned = strings.TrimSpace(cleaned)
	cleaned = regexp.MustCompile(`\s+`).ReplaceAllString(cleaned, " ")

	return cleaned
}

func main() {
	// load env, make sure you export api key
	apiKey := os.Getenv("YOUTUBE_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "ERROR: YOUTUBE_API_KEY is not set")
		os.Exit(1)
	}
	channelRes, err := http.Get("https://youtube.googleapis.com/youtube/v3/channels?part=contentDetails&forHandle=triplescosmos&key=" + apiKey + "&maxResults=5")

	// if invalid request
	if err != nil {
		panic(err)
	}
	defer channelRes.Body.Close()
	if channelRes.StatusCode != 200 {
		fmt.Println("Youtube Channels API not Available")
	}

	// get response body
	channelsBody, err := io.ReadAll(channelRes.Body)
	if err != nil {
		panic(err)
	}

	// unmarshal channels
	var channel ChannelItems
	json.Unmarshal(channelsBody, &channel)
	uploadsPlaylistId := channel.Items[0].ContentDetails.RelatedPlaylists.Uploads

	// get videos from playlist
	playlistRes, err := http.Get("https://youtube.googleapis.com/youtube/v3/playlistItems?part=snippet&key=" + apiKey + "&playlistId=" + uploadsPlaylistId + "&maxResults=1")

	// if invalid request
	if err != nil {
		panic(err)
	}
	defer playlistRes.Body.Close()
	if playlistRes.StatusCode != 200 {
		fmt.Println("Youtube PlaylistItems API not Available")
	}

	// get body
	playlistItemsBody, err := io.ReadAll(playlistRes.Body)
	if err != nil {
		panic(err)
	}

	// ummarshal
	var playlistItems PlaylistItems

	json.Unmarshal(playlistItemsBody, &playlistItems)
	videoTitle := playlistItems.Items[0].Snippet.Title
	videoId := playlistItems.Items[0].Snippet.ResourceId.VideoId

	fmt.Println(cleanTitle(videoTitle))
	fmt.Printf("Enjoy Your Content! https://www.youtube.com/watch?v=%s\n", videoId)
}
