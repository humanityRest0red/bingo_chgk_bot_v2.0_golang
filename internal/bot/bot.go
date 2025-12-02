package bot

import (
	"log"
	"os"

	"bingo-chgk-bot-v2.0-golang/internal"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

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
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	defer file.Close()
	log.SetOutput(file)

	bot, err := tgbotapi.NewBotAPI(internal.BotToken)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Авторизован как %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		log.Println(err)
	}

	return bot, updates
}
