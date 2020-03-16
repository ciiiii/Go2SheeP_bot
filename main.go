package main

import (
	"fmt"
	"os"

	"github.com/ciiiii/Go2SheeP_bot/pkg/translate"
	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {

	var (
		port           = os.Getenv("PORT")
		publicURL      = os.Getenv("PUBLIC_URL")
		token          = os.Getenv("TOKEN")
		translateAppId = os.Getenv("TRANSLATE_APP_ID")
		translateKey   = os.Getenv("TRANSLATE_KEY")
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
	tr := translate.NewTranslator(translateAppId, translateKey)
	bot.Handle(tb.OnText, func(m *tb.Message) {
		fmt.Printf("%s: %s\n", m.Sender.Username, m.Text)
		bot.Send(m.Sender, m.Text)
	})

	bot.Handle(tb.OnPhoto, func(m *tb.Message) {
		bot.Send(m.Sender, m.Photo.FileURL)
	})

	bot.Handle(tb.OnDocument, func(m *tb.Message) {
		bot.Send(m.Sender, m.Photo.FileURL)
	})

	bot.Handle("/en", func(m *tb.Message) {
		if len(m.Payload) == 0 {
			bot.Send(m.Sender, "⚠️input string is invalid.")
			return
		}
		result, err := tr.Translate("auto", "en", m.Payload)
		if err != nil {
			bot.Send(m.Sender, "👎translate failed.")
			return
		}
		bot.Send(m.Sender, fmt.Sprintf("👌:%s", result))
	})

	bot.Start()
}
