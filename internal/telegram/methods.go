package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
)

func (bot *TelegramBot) sendJSON(endpoint string, payload any) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("UNABLE TO MARSHAL JSON : %w", err)
	}

	reqURL := bot.API_URL + endpoint
	response, err := bot.Client.Post(reqURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("UNABLE TO SEND REQUEST : %w", err)
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

func (bot *TelegramBot) sendPhoto(chatId int64, imageBytes []byte, caption string) error {
	var requestBody bytes.Buffer

	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("photo", "image.jpg")
	if err != nil {
		return fmt.Errorf("UNABLE TO CREATE FORM FILE : %w", err)
	}

	if _, err := part.Write(imageBytes); err != nil {
		return fmt.Errorf("UNABLE TO WRITE IMAGE BYTES : %w", err)
	}

	if err := writer.WriteField("chat_id", strconv.FormatInt(chatId, 10)); err != nil {
		return fmt.Errorf("UNABLE TO WRITE CHAT ID : %w", err)
	}

	if caption != "" {
		if err := writer.WriteField("caption", caption); err != nil {
			return fmt.Errorf("UNABLE TO WRITE CAPTION : %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("UNABLE TO CLOSE WRITER : %w", err)
	}

	endpoint := bot.API_URL + "sendPhoto"
	req, err := http.NewRequest("POST", endpoint, &requestBody)
	if err != nil {
		return fmt.Errorf("UNABLE TO CREATE REQUEST : %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := bot.Client.Do(req)
	if err != nil {
		return fmt.Errorf("UNABLE TO SEND REQUEST : %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("TELEGRAM API Error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}
