package internal

import (
	"context"
	"time"
)

type rateLimitedYouTubeClient struct {
	next    YouTubeClient
	minWait time.Duration
	last    time.Time
}

// constructor 
func NewRateLimitedYouTubeClient(next YouTubeClient, minWait time.Duration) YouTubeClient {
	return &rateLimitedYouTubeClient{
		next:    next,
		minWait: minWait,
	}
}

// rate limit helper
func (r *rateLimitedYouTubeClient) wait() {
	now := time.Now()
	elapsed := now.Sub(r.last)

	if elapsed < r.minWait {
		time.Sleep(r.minWait - elapsed)
	}

	r.last = time.Now()
}

// wrapp calls of YoutubeClient
func (r * rateLimitedYouTubeClient) GetUploadsPlaylistId(ctx context.Context, channel string) (string, error) {
	r.wait()
	return r.next.GetUploadsPlaylistId(ctx, channel)
}

func (r * rateLimitedYouTubeClient) GetPlaylistItems(ctx context.Context, playListId string) (Video, error) {
	r.wait()
	return r.next.GetPlaylistItems(ctx, playListId)
}
