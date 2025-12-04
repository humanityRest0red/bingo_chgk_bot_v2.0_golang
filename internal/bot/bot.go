package bot

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"bingo-chgk-bot-v2.0-golang/internal"
	"bingo-chgk-bot-v2.0-golang/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var Articles = mapArticles()

func mapArticles() map[int]models.Article {
	var articles, _ = models.GetArticles()

	result := make(map[int]models.Article, len(articles))
	for i, a := range articles {
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

func setLogFile() {
	logDir := "../../logs"

	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	logFilePath := filepath.Join(logDir, "app.log")

	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	// defer file.Close() // no need to lof file

	log.SetOutput(file)
}

func botInitMust() (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel) {
	setLogFile()

	bot, err := tgbotapi.NewBotAPI(internal.BotToken)
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
