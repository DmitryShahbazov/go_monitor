package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"time"
)

func InitializeBot(BotTGToken string) tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(BotTGToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	return *bot
}

var (
	lastMessageTime time.Time
	lastMessageID   int64
)

func SendMessage(bot *tgbotapi.BotAPI, chatID int64, message string) {
	currentTime := time.Now()

	if currentTime.Sub(lastMessageTime) < time.Hour {
		editMsg := tgbotapi.NewEditMessageText(chatID, int(lastMessageID), message)
		editMsg.ParseMode = "HTML"
		editMsg.DisableWebPagePreview = true
		_, err := bot.Send(editMsg)
		if err.Error() == "Bad Request: MESSAGE_ID_INVALID" {
			msg := tgbotapi.NewMessage(chatID, message)
			msg.ParseMode = "HTML"
			msg.DisableWebPagePreview = true

			sentMsg, _ := bot.Send(msg)
			lastMessageTime = currentTime
			lastMessageID = int64(sentMsg.MessageID)
		}

	} else {
		msg := tgbotapi.NewMessage(chatID, message)
		msg.ParseMode = "HTML"
		msg.DisableWebPagePreview = true

		sentMsg, err := bot.Send(msg)
		if err == nil {
			deleteMsgConfig := tgbotapi.DeleteMessageConfig{
				ChatID:    chatID,
				MessageID: int(lastMessageID),
			}

			_, err = bot.Request(deleteMsgConfig)
			lastMessageTime = currentTime
			lastMessageID = int64(sentMsg.MessageID)
		}
	}
}
