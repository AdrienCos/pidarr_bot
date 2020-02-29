package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var keyboards []tb.InlineKeyboardMarkup = []tb.InlineKeyboardMarkup{}
var requestNb int64 = 0
var radarrToken string
var radarrHost string

const qualityProfile = 4 // HD 1080p config
const radarrPath string = "/movies"

// const radarrHost string = "localhost:7878"

type movieData struct {
	Year      int64        `json:"year"`
	Title     string       `json:"title"`
	ID        int64        `json:"tmdbId"`
	Quality   int          `json:"qualityProfileId"`
	TitleSlug string       `json:"titleSlug"`
	Images    []image      `json:"images"`
	Path      string       `json:"rootFolderPath"`
	Monitored bool         `json:"monitored"`
	Options   movieOptions `json:"addOptions"`
}

type image struct {
	CoverType string `json:"coverType"`
	Url       string `json:"url"`
}

type movieOptions struct {
	Search bool `json:"searchForMovie"`
}

func main() {
	// Grab the bot token and chat ID
	token := os.Getenv("PIDARR_TOKEN")
	if token == "" {
		log.Fatal("Unable to find the Telegram token")
	}
	chatID := os.Getenv("PIDARR_CHATID")
	if chatID == "" {
		log.Fatal("Unable to find the Telegram chat ID")
	}

	// Grab the Radarr API token
	radarrToken = os.Getenv("RADARR_TOKEN")
	if radarrToken == "" {
		log.Fatal("Unable to find the Radarr API token")
	}

	// Grab the Radarr host/port
	radarrHost = os.Getenv("RADARR_HOST")
	if radarrHost == "" {
		log.Fatal("Unable to find the Radarr host")
	}

	// Create a new bot
	b, err := tb.NewBot(tb.Settings{
		Token:    token,
		Poller:   &tb.LongPoller{Timeout: 10 * time.Second},
		Reporter: func(e error) { log.Print(e) },
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Print(("Bot created"))

	// Configure the bot's endpoints
	b.Handle("/movies", func(m *tb.Message) {
		SearchMovie(b, m)
	})
	b.Handle(tb.OnCallback, func(c *tb.Callback) {
		CallbackSearchMovie(b, c)
	})
	log.Print("Bot starting")
	// Start the bot
	b.Start()
}

func SearchMovie(b *tb.Bot, m *tb.Message) {
	if !IsCorrectChatID(m.Chat.ID) {
		log.Print("Bot called from ivalid chat")
		return
	}

	// Get the search term from the message
	search := m.Payload
	if search == "" {
		log.Print("/movies called without movie name")
		b.Send(m.Sender, "No movie title entered.")
		return
	}
	log.Printf("/movies called with search %s", search)
	// Search for the movie through the Radarr API
	values := make(map[string]string)
	values["apikey"] = radarrToken
	values["term"] = search
	u := NewUrl(radarrHost, "api/movie/lookup", values)
	response, err := http.Get(u.String())
	if err != nil {
		log.Print(err)
		b.Send(m.Sender, "Failed to search for movies, try again.")
		return
	}
	data, _ := ioutil.ReadAll(response.Body)
	// For all movies returned, extract the ID, name and year
	movies := []movieData{}
	json.Unmarshal(data, &movies)
	log.Printf("Found %d movies matching search", len(movies))
	if len(movies) == 0 {
		b.Send(m.Sender, "No movie found for this search.")
		return
	}
	// Create the new inline keyboard
	keyboard := tb.InlineKeyboardMarkup{}
	// For each movie, create a new button
	for _, movie := range movies {
		newButton := NewMovieButton(movie)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []tb.InlineButton{newButton})
	}
	// Return the list of movies to the User
	b.Send(m.Sender, "Here is the list of movies:", &tb.ReplyMarkup{InlineKeyboard: keyboard.InlineKeyboard})
	requestNb++
}

func CallbackSearchMovie(b *tb.Bot, c *tb.Callback) {
	// Send an empty response
	b.Respond(c, &tb.CallbackResponse{})
	// Extract the movie ID from c.Data
	// c.Data contains `requestNb|movieID`
	movieID := strings.Split(c.Data, "|")[1]
	log.Printf("Button for movie ID %s clicked", movieID)
	// Query the Radarr API for the movie information
	values := make(map[string]string)
	values["apikey"] = radarrToken
	values["tmdbId"] = movieID
	u := NewUrl(radarrHost, "api/movie/lookup/tmdb", values)
	response, err := http.Get(u.String())
	if err != nil {
		log.Print(err)
		b.Send(c.Sender, "Unable to find the movie, try again.")
		return
	}
	// Extract the movie data
	data, _ := ioutil.ReadAll(response.Body)
	movie := movieData{}
	json.Unmarshal(data, &movie)
	log.Printf("Found movie corresponding to ID: %s (%d)", movie.Title, movie.Year)
	movie.Monitored = true
	movie.Options.Search = true
	movie.Quality = qualityProfile
	movie.Path = radarrPath
	// Query the Radarr API to add and search for the movie
	postValues := make(map[string]string)
	postValues["apikey"] = radarrToken
	postUrl := NewUrl(radarrHost, "api/movie", postValues)
	body, _ := json.Marshal(movie)
	resp, err := http.Post(postUrl.String(), "application/json", bytes.NewBuffer((body)))
	// Check if the request went through
	if err != nil {
		log.Print(err)
		b.Send(c.Sender, "Unable to add the movie, try again.")
		return
	}
	// Check if the movie was successfully added
	code := resp.StatusCode
	if code != 201 {
		log.Print("Movie already tracked by Radarr")
		b.Send(c.Sender, "Movie already added to Radarr.")
	} else {
		log.Print("Movie added to Radarr")
		b.Send(c.Sender, "Movie added to Radarr, download will start soon.")
	}
}

func NewMovieButton(movie movieData) tb.InlineButton {
	movieName := movie.Title
	movieYear := movie.Year
	buttonText := fmt.Sprintf("%s (%d)", movieName, movieYear)
	movieID := movie.ID
	newButton := tb.InlineButton{
		Unique: strconv.FormatInt(requestNb, 10),
		Text:   buttonText,
		Data:   strconv.FormatInt(movieID, 10),
	}
	return newButton
}

func NewUrl(host string, path string, values map[string]string) url.URL {
	u := url.URL{
		Host:   host,
		Path:   path,
		Scheme: "http",
	}
	r := url.Values{}
	for key, value := range values {
		r.Add(key, value)
	}
	u.RawQuery = r.Encode()
	return u
}

// IsCorrentChatID returns whether the given chat ID is the same as the one set in PIDARR_CHATID
func IsCorrectChatID(c int64) bool {
	chatID, _ := strconv.ParseInt(os.Getenv("PIDARR_CHATID"), 10, 64)
	if chatID == c {
		return true
	} else {
		return false
	}
}
