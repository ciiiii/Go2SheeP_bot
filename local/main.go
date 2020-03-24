package main

import (
	"time"

	"github.com/ciiiii/Go2SheeP_bot/bot"
	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	ok := bot.ManageWebHook("delete")
	if !ok {
		panic("can't delete webhook")
	}
	p := &tb.LongPoller{Timeout: 15 * time.Second}
	b := bot.NewBot(p)
	b.Start()
}
