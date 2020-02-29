package main

import (
	"log"
	"time"

	"github.com/AdrienCos/pidarr_bot/internal/config"
	"github.com/AdrienCos/pidarr_bot/internal/endpoints"

	tb "gopkg.in/tucnak/telebot.v2"
)

var requestNb int64 = 0

func main() {
	// Create a new bot
	b, err := tb.NewBot(tb.Settings{
		Token:    config.Config.BotToken,
		Poller:   &tb.LongPoller{Timeout: 10 * time.Second},
		Reporter: func(e error) { log.Print(e) },
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Print(("Bot created"))

	// Configure the bot's endpoints
	b.Handle("/movies", func(m *tb.Message) {
		endpoints.SearchMovie(b, m, &requestNb)
	})
	b.Handle(tb.OnCallback, func(c *tb.Callback) {
		endpoints.CallbackSearchMovie(b, c)
	})

	// Start the bot
	log.Print("Bot starting")
	b.Start()
}
