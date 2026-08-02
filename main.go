package main

import (
	"fmt"
	"os"

	godotenv "github.com/joho/godotenv"
	logger "psghostelbot/internal/logger"
	tg "psghostelbot/internal/telegram"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	logger.ConfigLogger()

	tg.StartBot(os.Getenv("TELEGRAM_BOT_TOKEN"))
}
