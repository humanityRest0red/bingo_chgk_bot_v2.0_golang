package bot

import (
	"fmt"
	"math"
	"strings"

	"bingo-chgk-bot-v2.0-golang/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func handleCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	var text string
	var err error
	var msg tgbotapi.MessageConfig
	const maxLength = 4000

	switch update.Message.Command() {
	case "start", "help":
		text = helpText
		msg.ReplyMarkup = buildKeyboard()
	case "find":
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

	return err
}

func buildKeyboard() tgbotapi.ReplyKeyboardMarkup {
	buttonsTexts := []string{"Бинго", "Список статей", "Рандомная статья", "Статьи по темам"}
	var buttons []tgbotapi.KeyboardButton

	for _, text := range buttonsTexts {
		buttons = append(buttons, tgbotapi.NewKeyboardButton(text))
	}

	var rows [][]tgbotapi.KeyboardButton
	if len(buttons) > 0 {
		rows = append(rows, buttons)
	}

	return tgbotapi.NewReplyKeyboard(rows...)
}

func handleButtonPress(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	var response string
	switch update.Message.Text {
	case "Бинго":
		response = link("Бинго", config.BingoLink)
	case "Рандомная статья":
		response, _ = randomArticle()
	case "Статьи по темам":
		printArticles(bot, update)
		return nil
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	msg.ParseMode = tgbotapi.ModeMarkdown

	_, err := bot.Send(msg)
	return err
}

func printArticles(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	displayPage(bot, update, 1)
}

func displayPage(bot *tgbotapi.BotAPI, update tgbotapi.Update, pageNumber int) error {
	var err error
	const recordsPerPage = 30

	articles, _ := getArticles()
	articlesCount := len(articles)

	pagesCount := int(math.Ceil(float64(articlesCount)/float64(recordsPerPage))) + 1

	if pageNumber < 1 || pageNumber > pagesCount {
		return fmt.Errorf("выход за пределы страниц")
	}

	startIndex := articlesCount - 1 - (pageNumber-1)*recordsPerPage
	endIndex := articlesCount - 1 - pageNumber*recordsPerPage
	var text string
	for i := startIndex; i > endIndex && i > 0; i-- {
		text += fmt.Sprintf("%v. %s\n", articlesCount-i, link(articles[i].name, articles[i].link))
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
		// } else if strings.HasPrefix(callbackData, "type") {

		// 	key = call.data.split(SEP)[1]
		// 	articles, _ = getArticles()
		// 	for _, article := range articles {
		// 		if article.keys != "" {
		// 			keys = article.keys.split(',')
		// 			// if key in keys {
		// 			//     articles.append(article)
		// 			// }
		// 		}
		// 	}
		// 	var text string
		// 	articles.sort()
		// 	for i, article := range articles {
		// 		text += fmt.Sprintf("%d. %s\n", i+1, link(article.name, article.link))
		// 	}
		// 	// await call.message.answer(text=text, parse_mode='Markdown')
	}

	return nil
}

func selectTopics(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	topics, err := getTopics()
	if err != nil {
		return err
	}
	var buttons []tgbotapi.InlineKeyboardButton
	for _, topic := range topics {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(topic.name, fmt.Sprintf("type,%v", topic.key)))
	}

	var rows [][]tgbotapi.InlineKeyboardButton
	if len(buttons) > 0 {
		rows = append(rows, buttons)
	}

	markup := tgbotapi.InlineKeyboardMarkup{InlineKeyboard: rows}

	text := "Выберите тему:"
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.ReplyMarkup = markup
	if _, err := bot.Send(msg); err != nil {
		return fmt.Errorf("ошибка при отправке сообщения: %v", err)
	}

	return nil
}

// @dp.message_handler(commands=['find'])
// async def find_article(message: types.Message, command: Command.CommandObj):
//     if command.args:
//         text = ""
//         articles = get_articles()
//         for article in articles:
//             if command.args.lower() in article.name.lower():
//                 text += f"{link(article.name, article.link)}\n"
//         if text:
//             await message.answer(text, parse_mode='Markdown')
//         else:
//             await message.answer("По вашему запросу ничего не найдено")
//     else:
//         await message.answer("Укажите выражение после команды")
