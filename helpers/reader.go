package helpers

import (
	"bufio"
	"log"
	"os"
	"strings"
	"sync"
	"youtube_monitor/scrapers"
)

func ReadFile(filename string, wg *sync.WaitGroup, results chan<- scrapers.LiveCheckResult, scraper func(string, string, chan<- scrapers.LiveCheckResult)) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Ошибка при открытии файла: %s", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			log.Printf("Некорректный формат строки: %s", line)
			continue
		}
		url, nickname := parts[0], parts[1]

		wg.Add(1)
		go func(url, nickname string) {
			defer wg.Done()
			scraper(url, nickname, results)
		}(url, nickname)
	}
}
