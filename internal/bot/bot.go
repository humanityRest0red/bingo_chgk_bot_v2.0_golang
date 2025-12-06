package bot

import (
	"fmt"
	"log"
	"sync"

	"bingo-chgk-bot-v2.0-golang/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var ArticlesSlice, _ = models.GetArticles()
var ArticlesMap = mapArticles()

var userResponses = make(map[int64]chan string)
var mu sync.Mutex

func mapArticles() map[int]models.Article {
	if len(ArticlesSlice) == 0 {
		log.Fatalln("empty articles")
	}
	result := make(map[int]models.Article, len(ArticlesSlice))
	for i, a := range ArticlesSlice {
		result[i] = a
	}

	return result
}

func BotRun() {
	bot, updates := botInitMust()
	var err error

	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}

		if update.CallbackQuery != nil {
			log.Printf("[%s] Callback: %s", update.CallbackQuery.From.UserName, update.CallbackQuery.Data)
			err = handleCallback(bot, update)
		} else {
			log.Printf("[%s] Message: %s", update.Message.From.UserName, update.Message.Text)

			// collectResponses(update)

			if update.Message.IsCommand() {
				err = handleCommand(bot, update)
			} else {
				err = handleButtonPress(bot, update)
			}
		}

		if err != nil {
			log.Println(err)
		}
	}
}

func botInitMust() (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	log.Println("Авторизован как", bot.Self.UserName)
	fmt.Println("Авторизован как", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		log.Println(err)
	}

	return bot, updates
}
