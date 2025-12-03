package models

import (
	"encoding/json"
	"io"
 "log"
	"math/rand/v2"
	"os"
 "path/filepath"
	"slices"
	"strconv"
	"strings"
)

type Article struct {
	Index       int
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Keys        []string `json:"keys"`
}

func (a *Article) Link() string {
	return a.Name + " /" + strconv.Itoa(a.Index)
}

func (a *Article) Full() string {
	return a.Name + "\n\n" + a.Description
}

// func (a *Article) Name() string {
// 	return a.name
// }

func FilteredArticles(key string) ([]Article, error) {
	articles, err := GetArticles()
	if err != nil {
		return nil, err
	}

	filteredArticles := []Article{}
	for _, article := range articles {
		if slices.Contains(article.Keys, key) {
			filteredArticles = append(filteredArticles, article)
		}
	}

	slices.SortFunc(filteredArticles, func(a, b Article) int {
		return strings.Compare(a.Name, b.Name)
	})

	return filteredArticles, nil
}

func RandomArticle() (Article, error) {
	articles, err := GetArticles()
	if err != nil {
		return Article{}, err
	}

	i := rand.IntN(len(articles))
	return articles[i], nil
}

func GetArticles() ([]Article, error) {
currentDir, err := os.Getwd()
if err != nil {
    log.Println(err)
}
log.Println("Текущая рабочая директория:", currentDir)

basePath := filepath.Join("..", "..", "data")
absPath, err := filepath.Abs(filepath.Join(basePath, "test.json"))
 log.Println(absPath)
	file, err := os.OpenFile(absPath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var articles []Article
 log.Println("before unmarshall")
	err = json.Unmarshal(data, &articles)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(articles); i++ {
		articles[i].Index = i + 1
	}
 log.Println(len(articles))
	return articles, nil
}
