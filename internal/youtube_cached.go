package internal

import (
	"context"
	"sync"
	"time"
)

type cachedYouTubeClient struct {
	next YouTubeClient

	mu sync.Mutex

	playlistCache map[string]string
	videoCache    map[string]cachedVideo
}

type cachedVideo struct {
	video     Video
	expiresAt time.Time
}

// Constructor
func NewCachedYoutubeClient(next YouTubeClient) *cachedYouTubeClient {
	return &cachedYouTubeClient{
		next:          next,
		playlistCache: make(map[string]string),
		videoCache:    make(map[string]cachedVideo),
	}
}


// get playlistId from cache or cache from api call
func (c *cachedYouTubeClient) GetUploadsPlaylistId(ctx context.Context, channel string) (string, error) {
	// lock so only one go routinecan access
	c.mu.Lock()

	// access cache
	if playlistId, ok := c.playlistCache[channel]; ok {
		c.mu.Unlock()
		return playlistId, nil
	}

	// unlock mutex to make api call
	c.mu.Unlock()

	// make new api call
	playListId, err := c.next.GetUploadsPlaylistId(ctx, channel)
	if err != nil {
		return "", err
	}

	// add to cache whilst locked mutex
	c.mu.Lock()
	c.playlistCache[channel] = playListId
	c.mu.Unlock()

	return playListId, nil
}

// get playlistItems from cache or cache from api call
func (c *cachedYouTubeClient) GetPlaylistItems(ctx context.Context, playlistId string) (Video, error) {
	now := time.Now()

	// access cache while locked
	c.mu.Lock()
	if entry, ok := c.videoCache[playlistId]; ok && now.Before(entry.expiresAt) {
		c.mu.Unlock()
		return entry.video, nil
	}

	// unlock to make api call
	c.mu.Unlock()

	// add to cache

	video, err := c.next.GetPlaylistItems(ctx, playlistId)
	if err != nil {
		return Video{}, err
	}

	// lock mutex and add to cache
	c.mu.Lock()
	c.videoCache[playlistId] = cachedVideo{
		video:     video,
		expiresAt: now.Add(30 * time.Second),
	}
	c.mu.Unlock()

	return video, nil
}

