package main

import (
	"fmt"
	"github.com/ilya-sokolov/redbyte_bot/common"
	"github.com/ilya-sokolov/redbyte_bot/talks"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var c *cron.Cron

const (
	testPattern      = "2/2 * * * * *"
	oneDay           = "@daily"
	basePattern      = "0 10-18 10,24 * *"
	basePatternAdd12 = "0,59 10-18 10,12,24 * *"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}
	botToken := os.Getenv("bot_token")
	groupId, _ := strconv.Atoi(os.Getenv("group_id"))

	c = cron.New(
		cron.WithParser(
			cron.NewParser(
				cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)))
	s, _ := cron.ParseStandard(basePatternAdd12)

	fmt.Println("SCHEDULER:", s.Next(s.Next(time.Now())))
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	_, err = c.AddFunc(basePatternAdd12, func() {
		fmt.Println("NEXT START CRON:", s.Next(s.Next(time.Now())))
		msg := tgbotapi.NewMessage(int64(groupId), common.GetMessage())
		bot.Send(msg)
	})
	if err != nil {
		fmt.Println("Error: ", err)
	}
	c.Start()
	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		botSay(bot, update, int64(groupId))
	}
}

func botSay(bot *tgbotapi.BotAPI, update tgbotapi.Update, groupId int64) {
	if update.Message.Command() == "money" {
		msg := tgbotapi.NewMessage(groupId, "You have big money!")
		_, _ = bot.Send(msg)
		stopAndRestartCron(oneDay)
	}
	if len(update.Message.Text) != 0 {
		m := talks.NewMarkovChain("bot_dict.txt")
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(55)
		println(n)
		text := m.Generate(n)
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		message := strings.Split(update.Message.Text, " ")
		s0 := message[0]
		if s0 == "@SferaWoodpeckerBot" {
			msg := tgbotapi.NewMessage(groupId, text)
			_, _ = bot.Send(msg)
		} else if update.Message.Chat.ID != groupId {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
			_, _ = bot.Send(msg)
		}
	}

	if update.Message.Sticker != nil {
		msgSticker := tgbotapi.NewStickerShare(update.Message.Chat.ID, update.Message.Sticker.FileID)
		log.Printf("Sticker id = %s \n", update.Message.Sticker.FileID)
		msgSticker.ReplyToMessageID = update.Message.MessageID
		_, _ = bot.Send(msgSticker)
	}
}

func stopAndRestartCron(pattern string) {
	c.Stop()
	restart := cron.New()
	_, err := restart.AddFunc(pattern, func() {
		c.Start()
		restart.Stop()
	})
	if err != nil {
		fmt.Println("Error: ", err)
	}
	restart.Start()
}
