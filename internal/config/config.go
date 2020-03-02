package config

import (
	"log"
	"os"
	"strconv"
)

// Config is the global bot configuration
var Config Configuration

// Configuration is a struct that contains all the global configuration options
type Configuration struct {
	BotToken      string
	ChatID        int64
	RadarrToken   string
	RadarrHost    string
	RadarrPath    string
	RadarrQuality int64
}

const defaultRadarrHost string = "localhost:7878"
const defaultRadarrPath string = "/movies"
const defaultRadarrQuality string = "4" // HD 1080p

func init() {
	token := os.Getenv("PIDARR_TOKEN")
	if token == "" {
		log.Fatal("Unable to find the Telegram token")
	}
	Config.BotToken = token
	chatID, err := strconv.ParseInt(os.Getenv("PIDARR_CHATID"), 10, 64)
	if err != nil {
		log.Fatal("Unable to find the Telegram chat ID")
	}
	Config.ChatID = chatID

	// Grab the Radarr API token
	radarrToken := os.Getenv("RADARR_TOKEN")
	if radarrToken == "" {
		log.Fatal("Unable to find the Radarr API token")
	}
	Config.RadarrToken = radarrToken

	// Grab the Radarr host/port
	radarrHost := os.Getenv("RADARR_HOST")
	if radarrHost == "" {
		log.Print("Unable to find the Radarr host, defaulting to localhost:7878")
		radarrHost = defaultRadarrHost
	}
	Config.RadarrHost = radarrHost

	// Grab the Radarr path
	radarrPath := os.Getenv("RADARR_PATH")
	if radarrPath == "" {
		log.Print("Unable to find the Radarr path, defaulting to /movies")
		radarrPath = defaultRadarrPath
	}
	Config.RadarrPath = radarrPath

	// Grab the Radarr quality profile ID
	radarrQuality := os.Getenv("RADARR_QUALITY")
	if radarrQuality == "" {
		log.Print("Unable to find the Radarr quality, defaulting to 4 (HD 1080p)")
		radarrQuality = defaultRadarrQuality
	}
	Config.RadarrQuality, _ = strconv.ParseInt(radarrQuality, 10, 64)
}
