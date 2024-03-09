package settings

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
	"youtube_monitor/bot"
	"youtube_monitor/helpers"
	"youtube_monitor/scrapers"
)

func main() {
	resultsYt := make(chan scrapers.LiveCheckResult)
	resultsTw := make(chan scrapers.LiveCheckResult)
	resultsKc := make(chan scrapers.LiveCheckResult)

	var wg sync.WaitGroup
	var sb strings.Builder
	var resultsSliceYT []scrapers.LiveCheckResult
	var resultsSliceTW []scrapers.LiveCheckResult
	var resultsSliceKC []scrapers.LiveCheckResult

	for {
		resultsSliceYT = []scrapers.LiveCheckResult{}
		resultsSliceTW = []scrapers.LiveCheckResult{}
		resultsSliceKC = []scrapers.LiveCheckResult{}

		helpers.ReadFile("youtube_channels", &wg, resultsYt, scrapers.CheckIsLive)
		helpers.ReadFile("twitch_channels", &wg, resultsTw, scrapers.CheckIsLiveTwitch)
		helpers.ReadFile("kick_channels", &wg, resultsKc, scrapers.CheckIsLiveKick)

		go func() {
			for result := range resultsYt {
				resultsSliceYT = append(resultsSliceYT, result)
			}
			close(resultsYt)
		}()

		go func() {
			for result := range resultsTw {
				resultsSliceTW = append(resultsSliceTW, result)
			}
			close(resultsTw)
		}()

		go func() {
			for result := range resultsKc {
				resultsSliceKC = append(resultsSliceKC, result)
			}
			close(resultsKc)
		}()

		wg.Wait()
		sort.Sort(scrapers.ByViewCount(resultsSliceYT))
		sort.Sort(scrapers.ByViewCount(resultsSliceTW))
		sort.Sort(scrapers.ByViewCount(resultsSliceKC))

		if len(resultsSliceYT) > 0 {
			sb.WriteString("⭐️ <b>YOUTUBE</b> ⭐️\n\n")
			for _, result := range resultsSliceYT {
				//fmt.Printf("Nickname: %s, ViewCount: %d\n", result.Nickname, result.ViewCount)
				sb.WriteString(fmt.Sprintf("<b>%s [%d]</b>\n %s\n", result.Nickname, result.ViewCount, result.URL))
			}
			sb.WriteString("\n")
		}

		if len(resultsSliceTW) > 0 {
			sb.WriteString("⭐️ <b>TWITCH</b> ⭐️\n\n")
			for _, result := range resultsSliceTW {
				//fmt.Printf("Nickname: %s, ViewCount: %d\n", result.Nickname, result.ViewCount)
				sb.WriteString(fmt.Sprintf("<b>%s [%d]</b>\n %s\n", result.Nickname, result.ViewCount, result.URL))
			}
			sb.WriteString("\n")
		}

		if len(resultsSliceKC) > 0 {
			sb.WriteString("⭐️ <b>KICK</b> ⭐️\n\n")
			for _, result := range resultsSliceKC {
				//fmt.Printf("Nickname: %s, ViewCount: %d\n", result.Nickname, result.ViewCount)
				sb.WriteString(fmt.Sprintf("<b>%s [%d]</b>\n %s\n", result.Nickname, result.ViewCount, result.URL))
			}
		}

		if len(resultsSliceYT) > 0 || len(resultsSliceTW) > 0 || len(resultsSliceKC) > 0 {
			tgBot := bot.InitializeBot(BotTGToken)
			bot.SendMessage(&tgBot, ChannelID, sb.String())
			sb.Reset()
		}

		time.Sleep(SendTimeout * time.Minute)
	}

}
