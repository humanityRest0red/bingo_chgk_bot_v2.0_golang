package requests

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type MyService struct {
	baseUrl string
	limit   int
	client  http.Client
}

func NewMySerivce() *MyService {
	return &MyService{
		baseUrl: "https://gotquestions.online/search?search=",
		limit:   20,
		client:  http.Client{},
	}
}

func (ms *MyService) Test(substr string) (string, error) {
	url := ms.baseUrl + url.PathEscape(substr)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	body, err := ms.doRequestWithRetry(req)
	if err != nil {
		return "", fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}

	return string(body), nil
}

func (ms *MyService) doRequestWithRetry(req *http.Request) ([]byte, error) {
	for {
		resp, err := ms.client.Do(req)
		if err != nil {
			log.Println("Ошибка при выполнении запроса:", err)
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode == 429 {
			retryAfter := resp.Header.Get("Retry-After")
			waitSeconds := 0
			if retryAfter != "" {
				waitSeconds, _ = strconv.Atoi(retryAfter)
			} else {
				waitSeconds = 60
				waitSeconds = 1
			}
			log.Printf("Получен 429. Повтор через %d секунд...", waitSeconds)
			time.Sleep(time.Duration(waitSeconds) * time.Second)
			continue
		} else if resp.StatusCode != 200 {
			log.Println("Статус-код:", resp.Status)
		} else {
			if err != nil {
				log.Println("Ошибка при чтении тела ответа:", err)
			}
			return body, err
		}
	}
}
