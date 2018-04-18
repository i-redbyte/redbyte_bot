package main

import (
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"io/ioutil"
	"encoding/xml"
)

var config Config

type Config struct {
	TelegramBotToken string
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
		msg := tgbotapi.MessageConfig{}
		if update.Message.Text == "java" || update.Message.Text == "Java" {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, " Попробуй Go")
		} else {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, " Java, Kotlin, Go! We are doing good!")
			msg.ReplyToMessageID = update.Message.MessageID
		}
		bot.Send(msg)
	} else if update.Message.Sticker != nil {
		msg := tgbotapi.NewStickerShare(update.Message.Chat.ID, update.Message.Sticker.FileID)
		log.Printf("Sticker id = %s \n", update.Message.Sticker.FileID)
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)
	}
}
