package requests

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Question struct {
	ID      int    `json:"id"`
	Number  int    `json:"number"`
	Text    string `json:"text"`
	Answer  string `json:"answer"`
	Comment string `json:"comment"`
}

type QuestionWrapper struct {
	Questions []Question `json:"questions"`
}

func Find(str string) ([]Question, error) {
	ms := NewMySerivce()
	html, _ := ms.Test(str)
	html = strings.ReplaceAll(html, "\\\"", "\"")
	html = strings.ReplaceAll(html, "\\\\", "\\")

	begin := strings.Index(html, "\"questions\":")
	if begin == -1 {
		return nil, fmt.Errorf("тег 'questions' не найден")
	}
	html = html[begin:]

	end := strings.Index(html, ",\"count\":")
	if end == -1 {
		return nil, fmt.Errorf("тег 'count' не найден")
	}

	html = html[:end]

	html = "{" + html + "}"

	var questions QuestionWrapper
	err := json.Unmarshal([]byte(html), &questions)
	if err != nil {
		return nil, fmt.Errorf("requests json parser error: %v", err)
	}

	return questions.Questions, nil
}
