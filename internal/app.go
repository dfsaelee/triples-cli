package internal

import "context"

// Runs Business logic For http and domain types

type App struct {
	yt YouTubeClient
}

func NewApp(yt YouTubeClient) *App {
	return &App{yt: yt}
}

func (a *App) LatestVideo(ctx context.Context, channel string) (Video, error) {
	playlistId, err := a.yt.GetUploadsPlaylistId(ctx, channel)
	if err != nil {
		return Video{}, err
	}
	
	return a.yt.GetPlaylistItems(ctx, playlistId)
}