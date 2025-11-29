package bot

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const TABLE_NAME = "bingo.db"

type Record struct {
	name string
	link string
	keys string
}

func NewRecord(name, link, keys string) Record {
	return Record{
		name: name,
		link: link,
		keys: keys,
	}
}

// def __lt__(self, other):
//     return self.name < other.name

func getArticles() ([]Record, error) {
	var records []Record

	db, err := sql.Open("sqlite3", TABLE_NAME)
	if err != nil {
		return records, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM people")
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
