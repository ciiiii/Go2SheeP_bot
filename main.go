package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"fmt"
	"os"
)

func main() {
	var (
		port      = os.Getenv("PORT")
		publicURL = os.Getenv("PUBLIC_URL")
		token     = os.Getenv("TOKEN")
	)

	webhook := &tb.Webhook{
		Listen:   ":" + port,
		Endpoint: &tb.WebhookEndpoint{PublicURL: publicURL},
	}

	bot, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: webhook,
	})
	if err != nil {
		panic(err)
	}
	bot.Handle(tb.OnText, func(m *tb.Message) {
		fmt.Printf("%s: %s\n", m.Chat.Username, m.Text)
		bot.Send(m.Sender, m.Text)
	})

	bot.Handle(tb.OnPhoto, func(m *tb.Message) {
		bot.Send(m.Sender, m.Photo.FileURL)
	})

	bot.Handle(tb.OnDocument, func(m *tb.Message) {
		bot.Send(m.Sender, m.Photo.FileURL)
	})

	bot.Handle(tb.OnChannelPost, func(m *tb.Message) {

	})

	bot.Handle(tb.OnQuery, func(q *tb.Query) {

	})

	bot.Start()
}
