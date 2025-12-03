package models

// package main

// import (
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"
// 	"os"

// 	_ "github.com/mattn/go-sqlite3"
// )

// const TABLE_NAME = "../../data/bingo.db"

// func main() {
// 	articles, err := GetTopics()
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	jsonData, err := json.MarshalIndent(articles, "", "\t")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	os.WriteFile("topics.json", jsonData, 0644)
// 	for _, a := range articles {
// 		fmt.Println(a)
// 	}
// }

// type Topic struct {
// 	Name string `json:"name"`
// 	Key  string `json:"key"`
// }

// func GetTopics() ([]Topic, error) {
// 	var records []Topic

// 	db, err := sql.Open("sqlite3", TABLE_NAME)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer db.Close()

// 	rows, err := db.Query("SELECT * FROM Themes")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var record Topic
// 		err = rows.Scan(&record.Name, &record.Key)
// 		if err != nil {
// 			return nil, err
// 		}
// 		records = append(records, record)
// 	}

// 	err = rows.Err()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return records, nil
// }
