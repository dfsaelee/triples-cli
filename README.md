# triples-cli
Ever want to watch the latest TripleS content without moving away from your command line, in a quick, cached, API-limited way? This is for you... (cricket noises)

## Installation / Setup
Make sure you have Go installed

```
# Build
go build ./cmd/triples 

# Install to $GOPATH/bin or $GOBIN 
go install ./cmd/triples 
```
Verify Installation 
```
which triples
```

Set a YouTube data API v3 key as an environment variable:
```
# Bash / Zsh
export YOUTUBE_API_KEY=<YOUR_API_KEY>
```

Flags
## Command-Line Flags

| Flag | Description | Default |
|------|------------|---------|
| `-ch <channel-name>` | Specify the YouTube channel handle. | `triplescosmos` |
| `-health` | Run a quick health check to verify API key, network connectivity, and cache status. | N/A |

