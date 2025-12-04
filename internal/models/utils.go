package models

import (
	"fmt"
	"path/filepath"
)

const dataDir = "data"

var (
	ArticlesDataFilePath = filepath.Join(dataDir, "articles.json")
	TopicsDataFilePath   = filepath.Join(dataDir, "topics.json")
)

func Link(text, link string) string {
	return fmt.Sprintf("[%s](%s)", text, link)
}
