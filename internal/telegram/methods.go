package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (bot *TelegramBot) sendJSON(endpoint string, payload any) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	reqURL := bot.API_URL + endpoint
	response, err := bot.Client.Post(reqURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(response.Body)
		return fmt.Errorf("API Error %d: %s", response.StatusCode, string(bodyBytes))
	}

	return nil
}

func (bot *TelegramBot) sendMessage(message string, chatId *int64) error {
	err := bot.sendJSON("sendMessage", SendMessagePayload{
		ChatID:    *chatId,
		Text:      dedent(message),
		ParseMode: "Markdown",
	})

	return err
}

func (bot *TelegramBot) deleteMessage(chatId int64, messageId int) error {
	err := bot.sendJSON("deleteMessage", DeleteMessagePayload{
		ChatID:    chatId,
		MessageID: messageId,
	})

	return err
}
