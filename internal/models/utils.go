package models

import "path/filepath"

const dataDir = "data"

var (
	ArticlesDataFilePath = filepath.Join(dataDir, "articles.json")
	TopicsDataFilePath   = filepath.Join(dataDir, "topics.json")
)
