package main

type DBConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

type Access struct {
	UserID int64 `json:"user_id"`
}

type Config struct {
	TelegramToken string   `json:"telegram_token"`
	Database      DBConfig `json:"database"`
}
