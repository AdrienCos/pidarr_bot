package main

import (
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

type movieData struct {
	Year  int64  `json:"year"`
	Title string `json:"title"`
	ID    int64  `json:"tmdbId"`
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

	// Create a new bot
	b, err := tb.NewBot(tb.Settings{
		Token:    token,
		Poller:   &tb.LongPoller{Timeout: 10 * time.Second},
		Reporter: func(e error) { log.Print(e) },
	})
	if err != nil {
		log.Fatal(err)
	}

	// Configure the bot's endpoints
	b.Handle("/movies", func(m *tb.Message) {
		SearchMovie(b, m)
	})
	b.Handle(tb.OnCallback, func(c *tb.Callback) {
		fmt.Println("wtf")
		CallbackSearchMovie(b, c)
	})
	b.Handle("/test", func(m *tb.Message) {
		fmt.Println(m.Payload)
	})

	// Start the bot
	b.Start()
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

func SearchMovie(b *tb.Bot, m *tb.Message) {
	if IsCorrectChatID(m.Chat.ID) {
		// Get the search term from the message
		search := m.Payload
		if search == "" {
			log.Print("/movies called without movie name")
			b.Send(m.Sender, "No movie title entered")
			return
		}
		// Search for the movie through the Radarr API
		u := url.URL{
			Host:   "localhost:7878",
			Scheme: "http",
			Path:   "api/movie/lookup",
		}
		q := u.Query()
		q.Add("apikey", radarrToken)
		q.Add("term", search)
		u.RawQuery = q.Encode()
		response, err := http.Get(u.String())
		if err != nil {
			log.Fatal(err)
		}
		data, _ := ioutil.ReadAll(response.Body)
		// For all movies returned, extract the ID, name and year
		movies := []movieData{}
		json.Unmarshal(data, &movies)
		// Create the new inline keyboard
		keyboard := tb.InlineKeyboardMarkup{}
		// For each movie, create a new button
		for _, movie := range movies {
			fmt.Println(movie)
			movieName := movie.Title
			movieYear := movie.Year
			buttonText := fmt.Sprintf("%s (%d)", movieName, movieYear)
			movieID := movie.ID
			newButton := tb.InlineButton{
				Unique: strconv.FormatInt(requestNb, 10),
				Text:   buttonText,
				Data:   strconv.FormatInt(movieID, 10),
			}
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []tb.InlineButton{newButton})
		}
		// Return the list of movies to the User
		b.Send(m.Sender, "Here is the list of movies:", &tb.ReplyMarkup{InlineKeyboard: keyboard.InlineKeyboard})
		requestNb++
	}
}

func CallbackSearchMovie(b *tb.Bot, c *tb.Callback) {
	// c.Data contains `requestNb|movieID`
	fmt.Println(c.Data)
	// Send an empty response
	b.Respond(c, &tb.CallbackResponse{})
	// Extract the movie ID from c.Data
	movieID := strings.Split(c.Data, "|")[1]
	// Query the Radarr API for the movie information
	u := url.URL{
		Host:   "localhost:7878",
		Scheme: "http",
		Path:   "api/movie/lookup/tmdb",
	}
	r := u.Query()
	r.Add("apikey", radarrToken)
	r.Add("tmdbId", movieID)
	u.RawQuery = r.Encode()
	fmt.Println(u.String())
	response, err := http.Get(u.String())
	if err != nil {
		log.Print(err)
		b.Send(c.Sender, "Unable to find the movie, try again")
		return
	}
	data, _ := ioutil.ReadAll(response.Body)
	movie := movieData{}
	json.Unmarshal(data, &movie)
	fmt.Println(movie)
	// Query the Radarr API to add and search for the movie
	// TODO: fix my shit here
	b.Send(c.Sender, "hello")
}
