package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type YouTubeHTTPClient struct {
	apiKey string
	client *http.Client
}

func NewYoutubeHTTPClient(apiKey string, client *http.Client) *YouTubeHTTPClient {
	return &YouTubeHTTPClient{apiKey: apiKey, client: client}
}

type channelItemsResponse struct {
	Items []struct {
		ContentDetails struct {
			RelatedPlaylists struct {
				Uploads string `json:"uploads"`
			} `json:"relatedPlaylists"`
		} `json:"contentDetails"`
	} `json:"items"`
}

type playlistItemsResponse struct {
	Items []struct {
		Snippet struct {
			Title      string `json:"title"`
			ResourceId struct {
				VideoId string `json:"videoId"`
			} `json:"resourceId"`
		} `json:"snippet"`
	} `json:"items"`
}

// Get the playlist id
// modifying the YoutubeHTTPClient Function
func (c *YouTubeHTTPClient) GetUploadsPlaylistId(ctx context.Context, handle string) (string, error) {
	url := fmt.Sprintf(
		"https://youtube.googleapis.com/youtube/v3/channels?part=contentDetails&forHandle=%s&key=%s&maxResults=5",
		handle,
		c.apiKey,
	)
	channelRes, err := http.Get(url)

	// if invalid request
	if err != nil {
		return "", err
	}
	defer channelRes.Body.Close()
	if channelRes.StatusCode != 200 {
		return "", fmt.Errorf("YouTube Channels API returned status %d", channelRes.StatusCode)
	}

	// get response body
	channelsBody, err := io.ReadAll(channelRes.Body)
	if err != nil {
		return "", err
	}

	// unmarshal channels
	var channel channelItemsResponse
	if err := json.Unmarshal(channelsBody, &channel); err != nil {
		return "", err
	}
	return channel.Items[0].ContentDetails.RelatedPlaylists.Uploads, nil
}

func (c *YouTubeHTTPClient) GetPlaylistItems(ctx context.Context, playlistId string, maxResults int) ([]Video, error) {
	url := fmt.Sprintf(
		"https://youtube.googleapis.com/youtube/v3/playlistItems?part=snippet&key=%s&playlistId=%s&maxResults=%s",
		c.apiKey,
		playlistId,
		strconv.Itoa(maxResults),
	)
	playlistRes, err := http.Get(url)

	videos := make([]Video, 0, maxResults)

	if err != nil {
		return nil, err
	}

	defer playlistRes.Body.Close()
	if playlistRes.StatusCode != 200 {
		fmt.Println("Youtube PlaylistItems API not Available")
	}

	// get body
	playlistItemsBody, err := io.ReadAll(playlistRes.Body)
	if err != nil {
		return nil, err
	}

	// ummarshal
	var playlistItems playlistItemsResponse
	if err := json.Unmarshal(playlistItemsBody, &playlistItems); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if len(playlistItems.Items) < maxResults {
		maxResults = len(playlistItems.Items) // avoid out-of-bounds
	}

	for videoIndex := 0; videoIndex < maxResults; videoIndex++ {
		title := playlistItems.Items[videoIndex].Snippet.Title
		videoId := playlistItems.Items[videoIndex].Snippet.ResourceId.VideoId
		title = cleanTitle(title)
		videos = append(videos, Video{title, videoId})
	}
	return videos, nil
}
