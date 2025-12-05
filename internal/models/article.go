package models

import (
	"encoding/json"
	"io"
	"math/rand/v2"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

type Article struct {
	Index        int
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Associations []string `json:"associations"`
	Keys         []string `json:"keys"`
}

func (a *Article) Link() string {
	return a.Name + " /" + strconv.Itoa(a.Index)
}

func (a *Article) Full() string {
	var buf string
	if len(a.Associations) > 0 {
		buf += "\n\n**Ассоциации:**\n"
		for _, v := range a.Associations {
			buf += "— " + v + "\n"
		}
	}
	if len(a.Keys) > 0 {
		topics, _ := GetTopics()
		buf += "\n\n**Разделы:**\n"
		topicNames := []string{}
		for _, key := range a.Keys {
			for _, topic := range topics {
				if key == topic.Key {
					topicNames = append(topicNames, topic.Name)
					break
				}
			}
		}
		buf += strings.Join(topicNames, ", ")
	}

	return a.Link() + "\n\n" + a.Description + buf
}

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

func FilteredByWordArticles(substr string) []Article {
	filteredArticles := []Article{}
	if substr != "" {
		substr = strings.ToLower(substr)
		articles, _ := GetArticles()
		for _, article := range articles {
			if strings.Contains(strings.ToLower(article.Name), substr) ||
				strings.Contains(strings.ToLower(article.Description), substr) {
				filteredArticles = append(filteredArticles, article)
			}
		}
	}

	slices.SortFunc(filteredArticles, func(a, b Article) int {
		return strings.Compare(a.Name, b.Name)
	})

	return filteredArticles
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
	path, err := filepath.Abs(ArticlesDataFilePath)
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var articles []Article
	err = json.Unmarshal(data, &articles)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(articles); i++ {
		articles[i].Index = i + 1
	}

	return articles, nil
}
