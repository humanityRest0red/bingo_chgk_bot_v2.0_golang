package bot

import (
	"fmt"
	"math/rand/v2"
)

func randomRecord() (string, error) {
	records, err := getArticles()
	if err != nil {
		return "", err
	}

	i := rand.IntN(len(records))
	return link(records[i].name, records[i].link), nil
}

func link(text, link string) string {
	return fmt.Sprintf("[%s](%s)", text, link)
}

/*

SEP = ","

TABLE = get_articles()
USERS_ANKI = []


@dp.message_handler(commands=['list'])
async def print_records(message: types.Message):
    log(message)
    await display_the_page(message, page_number=1)


async def display_the_page(call, page_number):
    global TABLE
    TABLE = get_articles()
    records = TABLE
    records_count = len(records)
    pages_count = records_count // 30 + 1

    if page_number < 1 or page_number > pages_count:
        return

    markup = InlineKeyboardMarkup().add(
        InlineKeyboardButton("◀", callback_data=f"change_page{SEP}{page_number - 1}"),
        InlineKeyboardButton(f"{page_number}/{pages_count}", callback_data="null"),
        InlineKeyboardButton("▶", callback_data=f"change_page{SEP}{page_number + 1}"),
        InlineKeyboardButton("⏮ В начало", callback_data=f"change_page{SEP}{1}"),
        InlineKeyboardButton("В конец ⏭", callback_data=f"change_page{SEP}{pages_count}"),
    )
    text = ""
    for i in range(records_count - 1 - (page_number - 1) * 30, records_count - 1 - page_number * 30, -1):
        if i < 0:
            break
        text += f"{records_count - i}. {link(records[i].name, records[i].link)}\n"

    if type(call) == types.CallbackQuery:
        await call.message.edit_text(text=text, parse_mode='Markdown', reply_markup=markup)
    else:
        await call.answer(text=text, parse_mode='Markdown', reply_markup=markup)


@dp.callback_query_handler(text_startswith="change_page")
async def change_page(call: types.CallbackQuery):
    await call.answer()
    page_number = int(call.data.split(SEP)[1])

    await display_the_page(call, page_number)




@dp.message_handler(commands=['find'])
async def find_record(message: types.Message, command: Command.CommandObj):
    log(message)
    if command.args:
        text = ""
        records = get_articles()
        for record in records:
            if command.args.lower() in record.name.lower():
                text += f"{link(record.name, record.link)}\n"
        if text:
            await message.answer(text, parse_mode='Markdown')
        else:
            await message.answer("По вашему запросу ничего не найдено")
    else:
        await message.answer("Укажите выражение после команды")


@dp.message_handler(commands=['themes'])
async def select_theme(message):
    log(message)
    await message.delete()

    themes = get_themes()

    markup = InlineKeyboardMarkup()
    buttons = []
    for theme in themes:
        buttons.append(InlineKeyboardButton(text=theme[0], callback_data=f"type{SEP}{theme[1]}"))
    markup.add(*buttons)

    text = "Выберите тему:"
    await message.answer(text=text, parse_mode='Markdown', reply_markup=markup)


@dp.callback_query_handler(text_startswith="type")
async def send_category_records(call: types.CallbackQuery):
    key = call.data.split(SEP)[1]
    all_articles = get_articles()
    articles = []
    for article in all_articles:
        if article.keys:
            keys = article.keys.split(',')
            if key in keys:
                articles.append(article)

    text = ''
    articles.sort()
    for i, record in enumerate(articles):
        text += f"{i + 1}. {link(record.name, record.link)}\n"

    await call.message.answer(text=text, parse_mode='Markdown')


@dp.message_handler()
async def read_message(message: types.Message):
    log(message)
    if message.text.lower() == "старт" or message.text.lower() == "помощь":
        await send_welcome(message)
*/
