package models

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Topic struct {
	name string
	key  string
}

func (t *Topic) Name() string {
	return t.name
}

func (t *Topic) Key() string {
	return t.key
}

func GetTopics() ([]Topic, error) {
	var records []Topic

	db, err := sql.Open("sqlite3", TABLE_NAME)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM Themes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var record Topic
		err = rows.Scan(&record.name, &record.key)
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
