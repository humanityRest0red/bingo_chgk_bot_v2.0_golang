package bot

import (
	"log"

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
	bot, err := tgbotapi.NewBotAPI(botToken)
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
