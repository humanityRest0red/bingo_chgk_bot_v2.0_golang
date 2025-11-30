package bot

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const TABLE_NAME = "../../config/bingo.db"

type Record struct {
	name string
	link string
	keys sql.NullString
}

// func NewRecord(name, link string, keys sql.NullString) Record {
// 	return Record{
// 		name: name,
// 		link: link,
// 		keys: keys,
// 	}
// }

// def __lt__(self, other):
//     return self.name < other.name

func getArticles() ([]Record, error) {
	var records []Record

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
		var record Record
		err = rows.Scan(&record.name, &record.link, &record.keys)
		if err != nil {
			return []Record{}, err
		}
		records = append(records, record)
	}

	err = rows.Err()
	if err != nil {
		return []Record{}, err
	}

	return records, nil
}

type Topic struct {
	name string
	key  string
}

// func NewTopic(name, key string) Topic {
// 	return Topic{
// 		name: name,
// 		key:  key,
// 	}
// }

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
