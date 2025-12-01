package bot

import (
	"fmt"
	"math/rand/v2"
)

const helpText = `Команды:
/help  - выводит это сообщение
/find {выражение} - поиск статьи по выражению`

const (
	pageChangePrefix = "changePage:"
	topicsPrefix     = "topic:"
)

func randomArticle() (string, error) {
	articles, err := getArticles()
	if err != nil {
		return "Ошибка при отправке рандомной статьи", err
	}

	i := rand.IntN(len(articles))
	return link(articles[i].name, articles[i].link), nil
}

func link(text, link string) string {
	return fmt.Sprintf("[%s](%s)", text, link)
}

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
