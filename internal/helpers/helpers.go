package helpers

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/AdrienCos/pidarr_bot/internal/config"
	"github.com/AdrienCos/pidarr_bot/internal/types"

	tb "gopkg.in/tucnak/telebot.v2"
)

func NewMovieButton(movie types.MovieData, requestNb int64) tb.InlineButton {
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
	if config.Config.ChatID == c {
		return true
	} else {
		return false
	}
}
