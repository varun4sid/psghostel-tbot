package main

import (
	// "encoding/json"
	"fmt"
	godotenv "github.com/joho/godotenv"
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		fmt.Println("Error creating cookie jar:", err)
		return
	}

	client := &http.Client{
		Jar:     jar,
		Timeout: 10 * time.Second,
	}

	bot_auth_url := "https://api.telegram.org/bot" + os.Getenv("TELEGRAM_BOT_TOKEN") + "/getMe"
	fmt.Println("Checking bot authentication at:", bot_auth_url)

	response, err := client.Get(bot_auth_url)
	if err != nil {
		fmt.Println("Error fetching bot info:", err)
		return
	}
	defer response.Body.Close()

	fmt.Println("Response:", response)

	if response.StatusCode != http.StatusOK {
		fmt.Println("Failed to fetch bot info. Status code:", response.StatusCode)
		return
	}

	fmt.Println("Bot is running and authenticated successfully.")
}
