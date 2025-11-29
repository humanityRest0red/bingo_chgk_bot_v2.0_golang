package bot

import (
	"log"

	"bingo-chgk-bot-v2.0-golang/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const helpText = `Команды:
/help  - выводит это сообщение
/bingo - ссылка на основной файл
/list  - последние загруженные статьи
/rand  - случайная статья
/find {выражение} - поиск статьи по выражению
/topics - статьи по темам`

func BotRun() {
	bot, updates := botInitMust()

	maxLength := 4000

	for update := range updates {
		if update.Message == nil {
			continue
		}
		// update.Message.Text

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			var text string
			var err error
			switch update.Message.Command() {
			case "start", "help":
				text = helpText
			case "bingo":
				text = link("Бинго", config.BingoLink)
			case "list":
				text = "в работе"
			case "rand":
				text, err = randomRecord()
			case "find":
				text = "в работе"
			case "topics":
				text = "в работе"
			}

			for len(text) > 0 {
				chunk := text
				if len(chunk) > maxLength {
					chunk = text[:maxLength]
					text = text[maxLength:]
				} else {
					text = ""
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, chunk)
				msg.ParseMode = tgbotapi.ModeMarkdown

				bot.Send(msg)
			}

			if err != nil {
				log.Println(err)
			}
		}
	}
}

func botInitMust() (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel) {
	bot, err := tgbotapi.NewBotAPI(config.TokenBot)
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
