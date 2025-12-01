package bot

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const TABLE_NAME = "../../config/bingo.db"

type Article struct {
	name string
	link string
	keys sql.NullString
}

func getArticles() ([]Article, error) {
	var records []Article

	db, err := sql.Open("sqlite3", TABLE_NAME)
	if err != nil {
		return records, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM Articles")
	if err != nil {
		return records, err
	}
	defer rows.Close()

	for rows.Next() {
		var record Article
		err = rows.Scan(&record.name, &record.link, &record.keys)
		if err != nil {
			return []Article{}, err
		}
		records = append(records, record)
	}

	err = rows.Err()
	if err != nil {
		return []Article{}, err
	}

	return records, nil
}

type Topic struct {
	name string
	key  string
}

func getTopics() ([]Topic, error) {
	var records []Topic

	db, err := sql.Open("sqlite3", TABLE_NAME)
	if err != nil {
		return records, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM Themes")
	if err != nil {
		return records, err
	}
	defer rows.Close()

	for rows.Next() {
		var record Topic
		err = rows.Scan(&record.name, &record.key)
		if err != nil {
			return []Topic{}, err
		}
		records = append(records, record)
	}

	err = rows.Err()
	if err != nil {
		return []Topic{}, err
	}

	return records, nil
}
