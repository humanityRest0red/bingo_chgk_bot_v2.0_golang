package bot

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"bingo-chgk-bot-v2.0-golang/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func handleCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	// var text string
	var err error
	// const maxLength = 4000

	switch update.Message.Command() {
	case "start", "help":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, helpText)
		msg.ReplyMarkup = buildKeyboard()
		_, err = bot.Send(msg)
		return err
	case "find":
		err = findArticle(bot, update)
	case "log":
		err = sendLog(bot, update)
	default:
		n, err := strconv.Atoi(update.Message.Command())
		if err != nil {
			return err
		}

		article, exists := ArticlesMap[n-1]
		if !exists {
			return fmt.Errorf("статьи с указаннным номером нет")
		}

		return sendArticle(bot, update, article)
	}

	return err
}

func sendArticle(bot *tgbotapi.BotAPI, update tgbotapi.Update, article models.Article) error {
	text := article.Full()

	markup := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{
				// tgbotapi.NewInlineKeyboardButtonData("Google", "google:"+article.Name),
				// tgbotapi.NewInlineKeyboardButtonData("Вики", "wiki:"+article.Name),
				tgbotapi.NewInlineKeyboardButtonData("Рандомный вопрос", "questions:"+article.Name),
			},
		},
	}

	searchLink := models.Link("Google", "https://www.google.com/search?hl=ru&q="+url.QueryEscape(strings.ReplaceAll(article.Name, " ", "+")))
	text += "\n\n" + searchLink

	wikiLink := models.Link("Wikipedia", "https://ru.wikipedia.org/wiki/"+url.QueryEscape(strings.ReplaceAll(article.Name, " ", "_")))
	text += "\n" + wikiLink

	buf := strings.NewReplacer("(", "", ")", "").Replace(article.Name)
	questionsLink := models.Link("Вопросы в базе", "https://gotquestions.online/search?search="+url.QueryEscape(buf))
	text += "\n\n" + questionsLink

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)

	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.ReplyMarkup = markup
	if _, err := bot.Send(msg); err != nil {
		return fmt.Errorf("ошибка при отправке сообщения: %v", err)
	}

	return nil
}

func buildKeyboard() tgbotapi.ReplyKeyboardMarkup {
	var (
		buttonsTexts = []string{ /*"Бинго", */ "Список статей", "Рандомная статья", "Статьи по темам"}
		cols         = 2
		rows         = len(buttonsTexts)/cols + 1
		buttons      = make([][]tgbotapi.KeyboardButton, rows)
	)

	for i, text := range buttonsTexts {
		ind := i % cols
		buttons[ind] = append(buttons[ind], tgbotapi.NewKeyboardButton(text))
	}

	return tgbotapi.NewReplyKeyboard(buttons...)
}

func handleButtonPress(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	var response string
	switch update.Message.Text {
	// case "Бинго":
	// response = models.Link("Бинго", bingoLink)
	case "Список статей":
		printArticles(bot, update)
	case "Рандомная статья":
		article, err := models.RandomArticle()
		if err == nil {
			sendArticle(bot, update, article)
		}
	case "Статьи по темам":
		selectTopics(bot, update)
		return nil
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	msg.ParseMode = tgbotapi.ModeMarkdown

	_, err := bot.Send(msg)
	return err
}

func printArticles(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	displayPage(bot, update, 1) // pageNumber = 1
}

func displayPage(bot *tgbotapi.BotAPI, update tgbotapi.Update, pageNumber int) error {
	var err error
	const recordsPerPage = 30

	articles, _ := models.GetArticles()
	articlesCount := len(articles)

	pagesCount := int(math.Ceil(float64(articlesCount)/float64(recordsPerPage))) + 1

	if pageNumber < 1 || pageNumber > pagesCount {
		return fmt.Errorf("выход за пределы страниц")
	}

	startIndex := articlesCount - 1 - (pageNumber-1)*recordsPerPage
	endIndex := max(articlesCount-1-pageNumber*recordsPerPage, 0)
	var text string
	for i := startIndex; i >= endIndex; i-- {
		text += fmt.Sprintf("%v. %s\n", articlesCount-i, articles[i].Link())
	}

	markup := buildInlineKeyboard(pageNumber, pagesCount)

	if update.CallbackQuery == nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
		msg.ParseMode = tgbotapi.ModeMarkdown
		msg.ReplyMarkup = markup
		if _, err := bot.Send(msg); err != nil {
			return fmt.Errorf("ошибка при отправке сообщения: %v", err)
		}
	} else {
		callbackQuery := update.CallbackQuery
		_, err = bot.AnswerCallbackQuery(tgbotapi.CallbackConfig{
			CallbackQueryID: callbackQuery.ID,
			Text:            "Ваш ответ",
			ShowAlert:       false,
		})
		if err != nil {
			err = fmt.Errorf("ошибка при ответе на callback: %v", err)
			return err
		}

		editedMsg := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, text)
		editedMsg.ParseMode = tgbotapi.ModeMarkdown

		editedMsg.ReplyMarkup = &markup
		bot.Send(editedMsg)
	}

	return nil
}

func buildInlineKeyboard(currentPage, totalPages int) tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{
				tgbotapi.NewInlineKeyboardButtonData("◀", createPageChangeCommand(currentPage-1)),
				tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d/%d", currentPage, totalPages), "null"),
				tgbotapi.NewInlineKeyboardButtonData("▶", createPageChangeCommand(currentPage+1)),
			},
			{
				tgbotapi.NewInlineKeyboardButtonData("⏮ В начало", createPageChangeCommand(1)),
				tgbotapi.NewInlineKeyboardButtonData("В конец ⏭", createPageChangeCommand(totalPages)),
			},
		},
	}

	return keyboard
}

func handleCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	callbackData := update.CallbackQuery.Data

	if strings.HasPrefix(callbackData, pageChangePrefix) {
		pageNumber, err := extractPageNumber(callbackData)
		if err != nil {
			return fmt.Errorf("ошибка при извлечении номера страницы: %v", err)
		}
		displayPage(bot, update, pageNumber)
		return nil
	}

	var text string
	if strings.HasPrefix(callbackData, topicsPrefix) {
		key := callbackData[len(topicsPrefix):]
		var filteredArticles, err = models.FilteredArticles(key)
		if err != nil {
			return err
		}
		for i, article := range filteredArticles {
			text += fmt.Sprintf("%d. %s\n", i+1, article.Link())
		}
	} else if strings.HasPrefix(callbackData, "questions:") {
		return sendMessage(bot, update.Message.Chat.ID, "В разработке")
	}

	if text == "" {
		return nil
	}

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown

	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}

func selectTopics(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	topics, err := models.GetTopics()
	if err != nil {
		return err
	}

	var (
		cols    = 3
		rows    = len(topics) / cols
		buttons = make([][]tgbotapi.InlineKeyboardButton, rows)
	)

	for i, topic := range topics {
		ind := i / cols
		buttons[ind] = append(buttons[ind], tgbotapi.NewInlineKeyboardButtonData(topic.Name, fmt.Sprintf("%s%v", topicsPrefix, topic.Key)))
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите тему:")
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{InlineKeyboard: buttons}

	_, err = bot.Send(msg)

	return err
}

func findArticle(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	command := update.Message.Command()
	textAfterCommand := strings.TrimSpace(strings.TrimPrefix(update.Message.Text, "/"+command))

	if textAfterCommand == "" {
		return sendMessage(bot, update.Message.Chat.ID, "Укажите выражение после команды")
	}

	filteredArticles := models.FilteredByWordArticles(textAfterCommand)
	switch len(filteredArticles) {
	case 0:
		return sendMessage(bot, update.Message.Chat.ID, "По вашему запросу ничего не найдено")
	case 1:
		return sendArticle(bot, update, filteredArticles[0])
	default:
		return sendMultipleArticles(bot, update, filteredArticles)
	}
}

func sendLog(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	if update.Message.Chat.ID != 1077924714 {
		return nil
	}

	file, err := os.Open(filepath.Join("logs", "app.log"))
	if err != nil {
		log.Printf("Error opening log file: %v", err)
		return err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading log file: %v\n", err)
		return err
	}

	const maxLines = 50
	if len(lines) > maxLines {
		lines = lines[len(lines)-maxLines:]
	}

	text := strings.Join(lines, "\n")

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	_, err = bot.Send(msg)

	return err
}

func sendMultipleArticles(bot *tgbotapi.BotAPI, update tgbotapi.Update, articles []models.Article) error {
	var builder strings.Builder
	for _, article := range articles {
		builder.WriteString(article.Link() + "\n")
	}
	return sendMessage(bot, update.Message.Chat.ID, builder.String())
}

func sendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := bot.Send(msg)
	return err
}
