package helpers

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/AdrienCos/pidarr_bot/internal/config"
	"github.com/AdrienCos/pidarr_bot/internal/types"

	tb "gopkg.in/tucnak/telebot.v2"
)

// NewMovieButton creates a new InlineButton with the given movie data and returns it
func NewMovieButton(movie types.MovieData, requestNb int64, owned bool) tb.InlineButton {
	movieName := movie.Title
	movieYear := movie.Year
	var buttonText string
	if owned {
		buttonText = fmt.Sprintf("%s (%d) - In Colletion", movieName, movieYear)
	} else {
		buttonText = fmt.Sprintf("%s (%d)", movieName, movieYear)
	}
	movieID := movie.ID
	newButton := tb.InlineButton{
		Unique: strconv.FormatInt(requestNb, 10),
		Text:   buttonText,
		Data:   strconv.FormatInt(movieID, 10),
	}
	return newButton
}

// NewURL creates a new complete URL (including key-value options) and returns it
func NewURL(host string, path string, values map[string]string) url.URL {
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

// IsCorrectChatID returns whether the given chat ID is the same as the one set in PIDARR_CHATID
func IsCorrectChatID(c int64) bool {
	if config.Config.ChatID == c {
		return true
	}
	return false
}

// MovieInCollection returns whether the given movies is in the given collection
func MovieInCollection(movie types.MovieData, collection []types.MovieData) bool {
	for _, c := range collection {
		if c.ID == movie.ID {
			return true
		}
	}
	return false
}
