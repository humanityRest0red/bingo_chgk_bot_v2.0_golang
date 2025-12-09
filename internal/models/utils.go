package models

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const dataDir = "data"

var (
	articlesDataFilePath = filepath.Join(dataDir, "articles.json")
	topicsDataFilePath   = filepath.Join(dataDir, "topics.json")
)

var (
	articlesFileURL = os.Getenv("ARTICLES_FILE_LINK")
	topicsFileURL   = os.Getenv("TOPICS_FILE_LINK")
)

func Link(text, link string) string {
	return fmt.Sprintf("[%s](%s)", text, link)
}

func getFile(filePath, url string) error {
	output, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: %s", response.Status)
	}
	_, err = io.Copy(output, response.Body)

	return err
}
