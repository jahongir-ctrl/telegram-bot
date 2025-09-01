package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	//"go/types"
	"log"
	"os"
	"path/filepath"
	"time"
)

var bot *tgbotapi.BotAPI

func StartTelegramBot(token string) {
	if token == "" {
		log.Fatal("Telegram token is empty")
	}
	var err error
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Failed to start BOT: %v", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 35
	updates := bot.GetUpdatesChan(u)

	log.Printf("Bot started as %s", bot.Self.UserName)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if !isUserAllowed(update.Message.From.ID) {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Access denied"))
			continue
		}

		handleMessage(update.Message.Chat.ID)
	}
}
func handleMessage(chatID int64) {
	yesterday := time.Now().AddDate(0, 0, -1)
	//_, err := GenerateDailyReport(db, yesterday, "reports")
	filename := fmt.Sprintf("report_%s.txt", yesterday.Format("2006-01-02"))
	filepath := filepath.Join("reports", filename)

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		bot.Send(tgbotapi.NewMessage(chatID,
			fmt.Sprintf("Report not found for "+yesterday.Format("2006-01-02"))))
		return
	}
	//doc := tgbotapi.NewDocument(chatID, tgbotapi.FilePath(filepath))
	//doc.Caption = fmt.Sprintf("Report for %s", yesterday.Format("2006-01-02"))
	//
	//if _, err := bot.Send(doc); err != nil {
	//	bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Failed to send report: %v", err)))
	//}
	data, err := os.ReadFile(filepath)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, "Ошибка чтения отчета"))
		return
	}

	// Отправляем как текст
	bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Отчет за %s:\n\n%s", yesterday.Format("2006-01-02"), string(data))))

}

func isUserAllowed(userID int64) bool {
	data, err := os.ReadFile("access.txt")
	if err != nil {
		return false
	}
	var allowed []Access
	if err := json.Unmarshal(data, &allowed); err != nil {
		return false
	}
	for _, u := range allowed {
		if u.UserID == userID {
			return true
		}
	}
	//for _, id := range allowed {
	//	if id == userID {
	//		return true
	//	}
	//}
	return false
}
