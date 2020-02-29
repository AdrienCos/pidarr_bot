package endpoints

import (
	"log"
	"strings"

	"github.com/AdrienCos/pidarr_bot/internal/api"
	"github.com/AdrienCos/pidarr_bot/internal/config"
	"github.com/AdrienCos/pidarr_bot/internal/helpers"

	tb "gopkg.in/tucnak/telebot.v2"
)

func SearchMovie(b *tb.Bot, m *tb.Message, requestNb *int64) {
	if !helpers.IsCorrectChatID(m.Chat.ID) {
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
	// Query the Radarr API for the movies matching search
	movies, err := api.RadarrLookupName(search)
	if err != nil {
		log.Print(err)
		b.Send(m.Sender, "Failed to search for movies, try again.")
		return
	}
	log.Printf("Found %d movies matching search", len(movies))
	if len(movies) == 0 {
		b.Send(m.Sender, "No movie found for this search.")
		return
	}
	// Create the new inline keyboard
	keyboard := tb.InlineKeyboardMarkup{}
	// For each movie, create a new button
	for _, movie := range movies {
		newButton := helpers.NewMovieButton(movie, *requestNb)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []tb.InlineButton{newButton})
	}
	// Return the list of movies to the User
	b.Send(m.Sender, "Here is the list of movies:", &tb.ReplyMarkup{InlineKeyboard: keyboard.InlineKeyboard})
	*requestNb++
}

func CallbackSearchMovie(b *tb.Bot, c *tb.Callback) {
	// Send an empty response
	b.Respond(c, &tb.CallbackResponse{})
	// Extract the movie ID from c.Data
	// c.Data contains `requestNb|movieID`
	movieID := strings.Split(c.Data, "|")[1]
	log.Printf("Button for movie ID %s clicked", movieID)
	// Query the Radarr API for the movie information
	movie, err := api.RadarrLookupID(movieID)
	if err != nil {
		log.Print(err)
		b.Send(c.Sender, "Unable to find the movie, try again.")
		return
	}
	log.Printf("Found movie corresponding to ID: %s (%d)", movie.Title, movie.Year)
	// Set all the required fields for the API query
	movie.Monitored = true
	movie.Options.Search = true
	movie.Quality = config.Config.RadarrQuality
	movie.Path = config.Config.RadarrPath
	// Query the Radarr API to add and search for the movie
	code, err := api.RadarrAdd(movie)
	if err != nil {
		log.Print(err)
		b.Send(c.Sender, "Unable to add the movie, try again.")
		return
	}
	// Check if the movie was successfully added
	if code != 201 {
		log.Print("Movie already tracked by Radarr")
		b.Send(c.Sender, "Movie already added to Radarr.")
	} else {
		log.Print("Movie added to Radarr")
		b.Send(c.Sender, "Movie added to Radarr, download will start soon.")
	}
}
