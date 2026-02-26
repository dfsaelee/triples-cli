package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/term"

	"github.com/dfsaelee/triples-cli/internal"
)

var ch = flag.String("ch", "triplescosmos", "Enter youtube channel handle.")
var health = flag.Bool("health", false, "Check health.")

func runHealthCheck() {
	// check env
	apiKey := os.Getenv("YOUTUBE_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "ERROR: YOUTUBE_API_KEY is not set")
		os.Exit(1)
	} else {
		fmt.Println("Youtube Data API v3 Key Present")
	}

	// check cache
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		fmt.Println("Cannot resolve cache directory:", err)
		os.Exit(1)
	}

	// check if writeable
	testDir := filepath.Join(cacheDir, "triples")
	if err := os.MkdirAll(testDir, 0o755); err != nil {
		fmt.Println("Cannot create cache directory", err)
		os.Exit(1)
	}

	// api reachibility
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(
		ctx,
		"GET",
		"https://youtube.googleapis.com/youtube/v3/",
		nil,
	)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Cannot reach Youtube Api ", err)
		os.Exit(1)
	}
	res.Body.Close()
	fmt.Println("Youtube API reachable")
	fmt.Println("Health Check Passed")
}

func readKey() (key rune, err error) {
	fd := os.Stdin.Fd()
	var buf [1]byte

	oldState, err := term.MakeRaw(int(fd))
	if err != nil {
		return 0, err
	}
	defer term.Restore(int(fd), oldState)

	_, err = os.Stdin.Read(buf[:])

	return rune(buf[0]), nil
}

func main() {
	flag.Parse()
	channelHandle := *ch
	if *health {
		runHealthCheck()
		return
	}
	// load env, make sure you export api key
	apiKey := os.Getenv("YOUTUBE_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "ERROR: YOUTUBE_API_KEY is not set")
		os.Exit(1)
	}
	httpClient := &http.Client{
		Timeout: 3 * time.Second,
	}

	base := internal.NewYoutubeHTTPClient(apiKey, httpClient) // base youtube client
	rateLimited := internal.NewRateLimitedYouTubeClient(base, 200*time.Millisecond)
	cacheFile := os.TempDir() + "/triples_cache.json"
	cached := internal.NewCachedYoutubeClient(rateLimited, cacheFile) // cache
	app := internal.NewApp(cached)

	// index that scrolls through videos
	var videoIndex = 0
	ctx := context.Background()

	for {
		video, err := app.LatestVideo(ctx, channelHandle, videoIndex)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("\r[%d] \"%s\" : https://www.youtube.com/watch?v=%s      \n",
			videoIndex, video.Title, video.VideoId)
		key, err := readKey()
		
		switch key {
		case 'j':
			videoIndex++
		case 'k':	
			if (videoIndex > 0) {
				videoIndex--
			} else {
				fmt.Println("Top of List Reached")
			}
		case 'q':
			return
		}
	}
}
