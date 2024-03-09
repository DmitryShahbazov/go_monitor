package scrapers

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// CheckIsLive теперь принимает канал для передачи результатов и не возвращает непосредственно значения.
func CheckIsLive(channel, nickname string, results chan<- LiveCheckResult) {
	var connection_url string
	if strings.HasPrefix(channel, "@") {
		connection_url = "https://www.youtube.com/"
	} else {
		connection_url = "https://www.youtube.com/channel/"
	}
	url := connection_url + channel + "/live"

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return
	}

	req.Header.Add("pgrade-insecure-requests", "1")                                                                                                                     // Пример добавления заголовка Accept
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")                     // Пример добавления заголовка User-Agent
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7") // Пример добавления заголовка User-Agent
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	re_video_id := regexp.MustCompile(`videoDetails":{"videoId":"([A-Za-z0-9_-]{11})"`)
	video_id_matches := re_video_id.FindStringSubmatch(string(body))

	re_live := regexp.MustCompile(`LIVE_STREAM_OFFLINE`)
	live_matches := re_live.FindStringSubmatch(string(body))

	re := regexp.MustCompile(`originalViewCount":"(\d+)"`)
	matches := re.FindStringSubmatch(string(body))

	if len(matches) > 1 && matches[1] != "1" && len(live_matches) == 0 {
		a, _ := strconv.ParseInt(matches[1], 10, 32)
		URL := "https://youtu.be/" + video_id_matches[1]
		results <- LiveCheckResult{int32(a), nickname, nil, URL}
		time.Sleep(100 * time.Millisecond)
	}

}
