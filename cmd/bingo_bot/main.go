package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"bingo-chgk-bot-v2.0-golang/internal/bot"
)

func main() {
	mustLogInit()
	bot.BotRun()
}

func mustLogInit() {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Fatalln("Не удалось загрузить часовой пояс: ", err)
	}

	time.Local = loc
	file, err := os.OpenFile(filepath.Join("logs", "app.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file:", err)
	}
	// defer file.Close() // no need to log file

	log.SetOutput(file)
}
