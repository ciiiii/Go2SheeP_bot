package main

import (
	"github.com/ciiiii/Go2SheeP_bot/bot"
	"github.com/ciiiii/Go2SheeP_bot/config"
	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	ok := bot.ManageWebHook("set")
	if !ok {
		panic("can't set webhook")
	}
	webHook := &tb.Webhook{
		Listen:   ":" + config.Parser().Deploy.Port,
		Endpoint: &tb.WebhookEndpoint{PublicURL: config.Parser().Bot.PublicUrl},
	}
	b := bot.NewBot(webHook)
	b.Start()
}
