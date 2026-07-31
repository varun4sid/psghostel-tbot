package main

import (
	"fmt"
	"os"

	godotenv "github.com/joho/godotenv"
	tg "psghostelbot/internal/telegram"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	tg.StartBot(os.Getenv("TELEGRAM_BOT_TOKEN"))
}
