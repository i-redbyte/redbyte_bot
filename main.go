package main

import (
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"io/ioutil"
	"encoding/xml"
	"encoding/json"
	"net/http"
	"net/url"
)

var config Config

type Config struct {
	TelegramBotToken string
}

type WikiSearchResults struct {
	ready   bool
	Query   string
	Results []Result
}

type Result struct {
	Name, Description, URL string
}

func init() {
	xmlFile, err := ioutil.ReadFile("config.xml")
	if err != nil {
		log.Fatal(err)
	}
	xml.Unmarshal(xmlFile, &config)
}

func main() {
	bot, err := tgbotapi.NewBotAPI(config.TelegramBotToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		botSay(bot, update)
	}
}

func botSay(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if len(update.Message.Text) != 0 {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		ms, _ := urlEncoded(update.Message.Text)
		request := "https://ru.wikipedia.org/w/api.php?action=opensearch&search=" + ms + "&limit=3&origin=*&format=json"

		message := wikiAPI(request)
		for _, val := range message {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, val)
			bot.Send(msg)
		}

	} else if update.Message.Sticker != nil {
		msgSticker := tgbotapi.NewStickerShare(update.Message.Chat.ID, update.Message.Sticker.FileID)
		log.Printf("Sticker id = %s \n", update.Message.Sticker.FileID)
		msgSticker.ReplyToMessageID = update.Message.MessageID
		bot.Send(msgSticker)
	}
}

func (result *WikiSearchResults) UnmarshalJSON(bs []byte) error {
	var array []interface{}
	if err := json.Unmarshal(bs, &array); err != nil {
		return err
	}
	result.Query = array[0].(string)
	for i := range array[1].([]interface{}) {
		result.Results = append(result.Results, Result{
			array[1].([]interface{})[i].(string),
			array[2].([]interface{})[i].(string),
			array[3].([]interface{})[i].(string),
		})
	}
	return nil
}

func wikiAPI(request string) (answer []string) {
	slice := make([]string, 3) //3 элемента
	if response, err := http.Get(request); err != nil {
		slice[0] = "Википедия не отвечает"
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		wsr := &WikiSearchResults{}
		if err = json.Unmarshal([]byte(contents), wsr); err != nil {
			slice[0] = "Что-то не так, попробуйте изменить свой вопрос"
		}

		if !wsr.ready {
			slice[0] = "Что-то не так, попробуйте изменить свой вопрос"
		}

		for i := range wsr.Results {
			slice[i] = wsr.Results[i].URL
		}
	}
	return slice
}

func urlEncoded(str string) (string, error) {
	u, err := url.Parse(str)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
