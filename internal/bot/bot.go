package bot

import (
	"fmt"
	"log"

	"bingo-chgk-bot-v2.0-golang/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// var userResponses = make(map[int64]chan string)

func BotRun() {
	bot, updates := botInitMust()
	var err error

	go models.UpdateArticles()

	for update := range updates {
		go func(update tgbotapi.Update) {
			if update.Message == nil && update.CallbackQuery == nil {
				return
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
		}(update)
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
