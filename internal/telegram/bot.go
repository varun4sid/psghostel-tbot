package telegram

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func createBot(token string) *TelegramBot {
	return &TelegramBot{
		Token:   token,
		API_URL: fmt.Sprintf("https://api.telegram.org/bot%s/", token),
		Client:  &http.Client{Timeout: 70 * time.Second},
	}
}

func StartBot(token string, db *sql.DB) {
	psg := createBot(token)
	offset := 0

	for {
		updatedURL := fmt.Sprintf("%sgetUpdates?timeout=60&offset=%d", psg.API_URL, offset)

		response, err := psg.Client.Get(updatedURL)
		if err != nil {
			log.Printf("Network error while fetching updates : %v Retrying in 3s...", err)
			time.Sleep(3 * time.Second)
			continue
		}

		var updateResp UpdateResponse
		if err := json.NewDecoder(response.Body).Decode(&updateResp); err != nil {
			log.Printf("JSON decode error : %v", err)
			response.Body.Close()
			continue
		}
		response.Body.Close()

		for _, update := range updateResp.Result {
			if update.UpdateID >= offset {
				offset = update.UpdateID + 1
			}

			if update.Message != nil && strings.HasPrefix(update.Message.Text, "/") {
				handleBotCommand(psg, update.Message, db)
			}
		}
	}
}
