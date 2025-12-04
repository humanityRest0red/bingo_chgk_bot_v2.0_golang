package main

import (
	"log"
	"os"
	"path/filepath"

	"bingo-chgk-bot-v2.0-golang/internal/bot"
)

func main() {
	setLogFile()
	bot.BotRun()
}

func setLogFile() {
	file, err := os.OpenFile(filepath.Join("logs", "app.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	// defer file.Close() // no need to lof file

	log.SetOutput(file)
}
