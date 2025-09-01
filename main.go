package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
)

func main() {
	// читаем JSON конфиг
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Не удалось открыть config.json:", err)
	}
	defer configFile.Close()

	var cfg Config
	if err := json.NewDecoder(configFile).Decode(&cfg); err != nil {
		log.Fatal("Ошибка разбора config.json:", err)
	}

	// подключаемся к БД
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer db.Close()

	// отчёт за вчера
	yesterday := time.Now().AddDate(0, 0, -1)

	filePath, err := GenerateDailyReport(db, yesterday, "reports")
	if err != nil {
		log.Fatal("Ошибка генерации отчета:", err)
	}

	fmt.Printf("Отчет успешно сформирован: %s\n", filePath)

	StartTelegramBot(cfg.TelegramToken)
}
