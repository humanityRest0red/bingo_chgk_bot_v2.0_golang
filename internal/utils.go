package internal

import (
	"fmt"
	"os"
)

const helpText = `Команды:
/help  - выводит это сообщение
/find {выражение} - поиск статьи по выражению`

const (
	PageChangePrefix = "changePage:"
	TopicsPrefix     = "topic:"
)

var (
	bingoLink = os.Getenv("BINGO_BOT_LINK")
	BotToken  = os.Getenv("BINGO_BOT_TOKEN")
)

func HelpText() string {
	return helpText
}

func BingoLink() string {
	return Link("Бинго", bingoLink)
}

func Link(text, link string) string {
	return fmt.Sprintf("[%s](%s)", text, link)
}

func CreatePageChangeCommand(pageNumber int) string {
	return fmt.Sprintf("%s%d", PageChangePrefix, pageNumber)
}

func ExtractPageNumber(callbackData string) (int, error) {
	var pageNumber int
	n, err := fmt.Sscanf(callbackData, "changePage:%d", &pageNumber)
	if err != nil {
		return 0, err
	}
	if n != 1 {
		return 0, fmt.Errorf("не удалось извлечь номер страницы из: %s", callbackData)
	}
	return pageNumber, nil
}

// func ExtractGoogle(callbackData string) (string, error) {
// 	var str string
// 	n, err := fmt.Sscanf(callbackData, "changePage:%d", &pageNumber)
// 	if err != nil {
// 		return 0, err
// 	}
// 	if n != 1 {
// 		return 0, fmt.Errorf("не удалось извлечь номер страницы из: %s", callbackData)
// 	}
// 	return pageNumber, nil
// }
