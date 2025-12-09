package models

import (
	"encoding/json"
	"io"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Topic struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

func GetTopics() ([]Topic, error) {
	err := getFile(topicsDataFilePath, topicsFileURL)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(topicsDataFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var articles []Topic
	err = json.Unmarshal(data, &articles)
	if err != nil {
		return nil, err
	}

	return articles, nil
}
