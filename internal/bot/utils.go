package bot

import (
	"fmt"
	"os"
)

const helpText = `Команды:
/help  - выводит это сообщение`

// /find {выражение} - поиск статьи по выражению`

const (
	pageChangePrefix = "changePage:"
	topicsPrefix     = "topic:"
)

var (
	// bingoLink = os.Getenv("BINGO_BOT_LINK")
	botToken = os.Getenv("BINGO_BOT_TOKEN")
)

func createPageChangeCommand(pageNumber int) string {
	return fmt.Sprintf("%s%d", pageChangePrefix, pageNumber)
}

func extractPageNumber(callbackData string) (int, error) {
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
