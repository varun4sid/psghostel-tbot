package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"psghostelbot/internal/db"
	"psghostelbot/internal/logger"
	tg "psghostelbot/internal/telegram"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	logger.ConfigLogger()
	database, err := db.NewConnection(os.Getenv("DB_URL"))
	if err != nil {
		slog.Error("Error connecting to database:", "error", err)
		os.Exit(1)
		return
	}
	defer database.Close()

	tg.StartBot(os.Getenv("TELEGRAM_BOT_TOKEN"), database)
}
