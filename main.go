package main

import (
	"fmt"
	"log/slog"
	"os"

	"psghostelbot/internal/db"
	"psghostelbot/internal/logger"
	tg "psghostelbot/internal/telegram"

	"github.com/joho/godotenv"
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

	psg := tg.CreateBot(os.Getenv("TELEGRAM_BOT_TOKEN"))
	scheduler := tg.CreateScheduler(database, psg)
	scheduler.StartAsync()

	tg.StartBot(psg, database)
}
