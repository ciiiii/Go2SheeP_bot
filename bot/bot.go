package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ciiiii/Go2SheeP_bot/config"
	"github.com/ciiiii/Go2SheeP_bot/pkg/cos"
	"github.com/ciiiii/Go2SheeP_bot/pkg/translate"
	"github.com/ciiiii/Go2SheeP_bot/pkg/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	cosService     *cos.Client
	storedMessages map[int]M
)

type M struct {
	Id     int
	ChatId int64
}

func NewBot(p tb.Poller) *tb.Bot {
	storedMessages = make(map[int]M, 0)
	bot, err := tb.NewBot(tb.Settings{
		Token:  config.Parser().Bot.Token,
		Poller: p,
	})
	if err != nil {
		panic(err)
	}
	tr := translate.NewTranslator(config.Parser().Translate.AppId, config.Parser().Translate.Key)
	cosService = cos.NewCos(config.Parser().Cos.Bucket, config.Parser().Cos.Region, config.Parser().Cos.SecretId, config.Parser().Cos.SecretKey)

	bot.Handle(tb.OnText, func(m *tb.Message) {
		bot.Reply(m, m.Text)
	})

	bot.Handle(tb.OnPhoto, func(m *tb.Message) {
		fileUrl, err := utils.GetFileUrl(config.Parser().Bot.Token, m.Photo.FileID)
		if err != nil {
			bot.Send(m.Sender, "👎 something wrong")
			fmt.Println(err)
			return
		}
		bot.Send(m.Sender, "💤 start downloading\n")
		tmpFile, err := utils.DownloadAsTmp(fileUrl)
		defer utils.DeleteFile(tmpFile)
		if err != nil {
			bot.Send(m.Sender, "👎 something wrong")
			fmt.Println(err)
			return
		}
		bot.Send(m.Sender, "💤 start uploading\n")
		cosUrl, err := cosService.Upload(tmpFile)
		if err != nil {
			bot.Send(m.Sender, "👎 something wrong")
			fmt.Println(err)
			return
		}
		bot.Send(m.Sender, "✅ "+cosUrl)
	})

	bot.Handle(tb.OnDocument, func(m *tb.Message) {
		bot.Send(m.Sender, m.Photo.FileURL)
	})

	bot.Handle("/en", func(m *tb.Message) {
		if len(m.Payload) == 0 {
			bot.Reply(m, "⚠️ input string is invalid")
			return
		}
		result, err := tr.Translate("auto", "en", m.Payload)
		if err != nil {
			bot.Reply(m, "👎 translate failed")
			return
		}
		bot.Reply(m, fmt.Sprintf("👌 %s", result))
	})

	bot.Handle("/cn", func(m *tb.Message) {
		if len(m.Payload) == 0 {
			bot.Reply(m, "⚠️ input string is invalid")
			return
		}
		result, err := tr.Translate("auto", "en", m.Payload)
		if err != nil {
			bot.Reply(m, "👎 translate failed")
			return
		}
		bot.Reply(m, fmt.Sprintf("👌 %s", result))
	})

	bot.Handle("/images", func(m *tb.Message) {
		bot.Send(m.Sender, "Manage Images", gen(true, ""))
	})

	bot.Handle(tb.OnCallback, func(callback *tb.Callback) {
		data := utils.CleanInvalidStr(callback.Data)
		prefix := data[:1]
		marker := data[1:]
		switch prefix {
		case "+":
			bot.Edit(callback.Message, "Manage Images", gen(false, marker))
			break
		case ".":
			go func() {
				if oldMessage, ok := storedMessages[callback.Sender.ID]; ok {
					message := &tb.Message{ID: oldMessage.Id, Chat: &tb.Chat{ID: oldMessage.ChatId}}
					bot.Delete(message)
				}
			}()
			photo := &tb.Photo{
				File:    tb.FromURL(fmt.Sprintf("http://%s.cos.%s.myqcloud.com/%s", config.Parser().Cos.Bucket, config.Parser().Cos.Region, marker)),
				Caption: "",
			}
			newMessage, err := bot.Send(callback.Sender, photo)
			if err != nil {
				fmt.Println(err)
				break
			}
			storedMessages[callback.Sender.ID] = M{
				Id:     newMessage.ID,
				ChatId: newMessage.Chat.ID,
			}
			break
		case "-":
			go func() {
				if oldMessage, ok := storedMessages[callback.Sender.ID]; ok {
					message := &tb.Message{ID: oldMessage.Id, Chat: &tb.Chat{ID: oldMessage.ChatId}}
					bot.Delete(message)
				}
			}()
			bot.Delete(callback.Message)
			break
		default:
			fmt.Println(prefix, marker)
		}
	})

	return bot
}

func gen(init bool, marker string) *tb.ReplyMarkup {
	imageList, _ := cosService.List(marker, 3)
	var replyKeyboard [][]tb.InlineButton
	if len(imageList) != 0 {
		if init {

		}
		album := make(tb.Album, len(imageList))
		keyboard := make([]tb.InlineButton, len(imageList))
		for i, image := range imageList {
			album[i] = &tb.Photo{File: tb.FromURL(image), Width: 1, Height: 1}
			keyboard[i] = tb.InlineButton{
				Unique: "." + image,
				Text:   image,
			}
		}
		nextButton := tb.InlineButton{
			Unique: "+" + imageList[len(imageList)-1],
			Text:   "next",
		}
		replyKeyboard = [][]tb.InlineButton{
			keyboard,
			{nextButton},
		}
	} else {
		nullButton := tb.InlineButton{
			Unique: "-",
			Text:   "exit",
		}
		replyKeyboard = [][]tb.InlineButton{
			{nullButton},
		}
	}

	return &tb.ReplyMarkup{
		InlineKeyboard: replyKeyboard,
	}
}

type result struct {
	Ok     bool `json:"ok"`
	Result bool `json:"result"`
}

func ManageWebHook(option string) bool {
	var url string
	switch option {
	case "set":
		url = fmt.Sprintf("https://api.telegram.org/bot%s/setwebhook?url=%s", config.Parser().Bot.Token, config.Parser().Bot.PublicUrl)
		break
	case "delete":
		url = fmt.Sprintf("https://api.telegram.org/bot%s/deletewebhook", config.Parser().Bot.Token)
		break
	default:
		return false
	}

	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	var r result
	if err != json.Unmarshal(body, &r) {
		return false
	}
	if r.Result && r.Ok {
		return true
	}
	return false
}
