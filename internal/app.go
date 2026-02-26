package internal

import (
	"context"
	"fmt"
)

// Runs Business logic For http and domain types

type App struct {
	yt YouTubeClient
}

func NewApp(yt YouTubeClient) *App {
	return &App{yt: yt}
}

func (a *App) LatestVideo(ctx context.Context, channel string, videoIndex int) (Video, error) {
	playlistId, err := a.yt.GetUploadsPlaylistId(ctx, channel)
	if err != nil {
		return Video{}, err
	}

	maxResults := ((videoIndex / 10) + 1) * 10

	videos, err := a.yt.GetPlaylistItems(ctx, playlistId, maxResults)
	if err != nil {
		return Video{}, err
	}

	if videoIndex >= len(videos) {
		return Video{}, fmt.Errorf("video index %d out of range (have %d)", videoIndex, len(videos))
	}

	if videoIndex >= 50 {
		return Video{}, fmt.Errorf("video index %d out of max range of api (have %d)", videoIndex, len(videos))
	}

	return videos[videoIndex], nil
}
