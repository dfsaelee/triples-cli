package internal

import (
    "context"
)

// Interface that defines and defines domain types
type Video struct {
    Title string
    VideoId string
}

// YoutubeHTTPClient satisfies interface with these methods
type YouTubeClient interface {
	GetUploadsPlaylistId(ctx context.Context, channel string) (string, error)
	GetPlaylistItems(ctx context.Context, playlistID string) (Video, error)
}