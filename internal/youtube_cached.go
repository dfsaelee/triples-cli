package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type cachedYouTubeClient struct {
	next YouTubeClient

	mu sync.Mutex

	playlistCache map[string]string
	videoCache    map[string]cachedVideos

	cacheFile string
}

type cachedVideos struct {
	Video     []Video   `json:"video"`
	ExpiresAt time.Time `json:"expires_at"`
}

type persistentCache struct {
	PlaylistCache map[string]string       `json:"playlist_cache"` // caches playlistId
	VideoCache    map[string]cachedVideos `json:"video_cache"`
}

// Constructor
func NewCachedYoutubeClient(next YouTubeClient, cacheFile string) YouTubeClient {
	c := &cachedYouTubeClient{
		next:          next,
		playlistCache: make(map[string]string),
		videoCache:    make(map[string]cachedVideos),
		cacheFile:     cacheFile,
	}
	c.loadCache()
	return c
}

// persistence
func (c *cachedYouTubeClient) loadCache() {
	data, err := os.ReadFile(c.cacheFile)
	if err != nil {
		return
	}
	var pc persistentCache
	if err := json.Unmarshal(data, &pc); err != nil {
		fmt.Println("cache parsing failed:", err)
		return
	}

	c.mu.Lock()
	c.playlistCache = pc.PlaylistCache
	c.videoCache = pc.VideoCache
	c.mu.Unlock()
}

func (c *cachedYouTubeClient) saveCache() {
	c.mu.Lock()
	pc := persistentCache{
		PlaylistCache: c.playlistCache,
		VideoCache:    c.videoCache,
	}
	c.mu.Unlock()

	data, _ := json.MarshalIndent(pc, "", " ")

	os.MkdirAll(filepath.Dir(c.cacheFile), 0o755)
	_ = os.WriteFile(c.cacheFile, data, 0o644)
}

// get playlistId from cache or cache from api call
func (c *cachedYouTubeClient) GetUploadsPlaylistId(ctx context.Context, channel string) (string, error) {
	// lock so only one go routinecan access
	c.mu.Lock()

	// access cache (cache already loaded from constructor)
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

	c.saveCache()
	return playListId, nil
}

// get playlistItems from cache or cache from api call
func (c *cachedYouTubeClient) GetPlaylistItems(ctx context.Context, playlistId string, maxResults int) ([]Video, error) {
	now := time.Now()

	// access cache while locked
	c.mu.Lock()
	if entry, ok := c.videoCache[playlistId]; ok &&
		now.Before(entry.ExpiresAt) &&
		maxResults <= len(entry.Video) {
		c.mu.Unlock()
		return entry.Video, nil
	}

	// unlock to make api call
	c.mu.Unlock()

	// add to cache

	videos, err := c.next.GetPlaylistItems(ctx, playlistId, maxResults)
	if err != nil {
		return nil, err
	}

	// lock mutex and add to cache
	c.mu.Lock()
	c.videoCache[playlistId] = cachedVideos{
		Video:     videos,
		ExpiresAt: now.Add(30 * time.Second),
	}
	c.mu.Unlock()

	c.saveCache()
	return videos, nil
}
