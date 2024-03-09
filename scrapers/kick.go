package scrapers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func CheckIsLiveKick(channel, nickname string, results chan<- LiveCheckResult) {
	url := "https://kick.com/api/v2/channels/" + channel + "/livestream" // Замените на нужный URL
	client := &http.Client{}                                             // Используйте http.DefaultClient, если нет специфических настроек

	// Создание нового запроса
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return
	}
	req2, err := http.NewRequest("GET", "https://kick.com/", nil)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return
	}

	// Добавление заголовков к запросу
	req.Header.Add("pgrade-insecure-requests", "1")                                                                                                                     // Пример добавления заголовка Accept
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")                     // Пример добавления заголовка User-Agent
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7") // Пример добавления заголовка User-Agent
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req2.Header.Add("pgrade-insecure-requests", "1")                                                                                                                     // Пример добавления заголовка Accept
	req2.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")                     // Пример добавления заголовка User-Agent
	req2.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7") // Пример добавления заголовка User-Agent
	req2.Header.Add("Accept-Language", "en-US,en;q=0.5")                                                                                                                 // Пример добавления заголовка User-Agent

	resp2, err := client.Do(req2)
	_, err = io.ReadAll(resp2.Body)
	defer resp2.Body.Close()

	req.Header = req2.Header
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при отправке запроса:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		results <- LiveCheckResult{0, "", err, ""}
		return
	}
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatal(err)
	}

	if dataMap, ok := data["data"].(map[string]interface{}); ok {
		if viewers, ok := dataMap["viewers"].(float64); ok {
			URL := "https://kick.com/" + channel
			results <- LiveCheckResult{int32(viewers), nickname, nil, URL}
		} else {
			fmt.Println("Ошибка: Не удается привести viewers к типу int")
		}
	}
}
