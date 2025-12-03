package models

import (
	"database/sql"
	"math/rand/v2"
	"slices"
	"strings"

	"bingo-chgk-bot-v2.0-golang/internal"
)

type Article struct {
	name string
	link string
	keys sql.NullString
}

func (a *Article) Link() string {
	return internal.Link(a.name, a.link)
}

func (a *Article) Name() string {
	return a.name
}

func FilteredArticles(key string) ([]Article, error) {
	articles, err := GetArticles()
	if err != nil {
		return nil, err
	}

	filteredArticles := []Article{}
	for _, article := range articles {
		if article.keys.Valid {
			keys := strings.Split(article.keys.String, ",")
			if slices.Contains(keys, key) {
				filteredArticles = append(filteredArticles, article)
			}
		}
	}

	slices.SortFunc(filteredArticles, func(a, b Article) int {
		return strings.Compare(a.name, b.name)
	})

	return filteredArticles, nil
}

func RandomArticle() (string, error) {
	articles, err := GetArticles()
	if err != nil {
		return "Ошибка при отправке рандомной статьи", err
	}

	i := rand.IntN(len(articles))
	return articles[i].Link(), nil
}

func GetArticles() ([]Article, error) {
	db, err := sql.Open("sqlite3", TABLE_NAME)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM Articles")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []Article
	var record Article
	for rows.Next() {
		err = rows.Scan(&record.name, &record.link, &record.keys)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return records, nil
}
