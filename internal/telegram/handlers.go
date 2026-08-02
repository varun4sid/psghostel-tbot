package telegram

import (
	"bytes"
	"database/sql"
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

func insertStudent(db *sql.DB, rollno, encryptedPassword string, chatID int64) (error, int) {
	existingChat := db.QueryRow("SELECT chat_id FROM students WHERE chat_id = $1", chatID)
	var existingChatID string
	err := existingChat.Scan(&existingChatID)
	if err != nil && err != sql.ErrNoRows {
		return err, 10
	}

	_, err = db.Exec(`
		INSERT INTO students (rollno, password, chat_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (chat_id) DO UPDATE SET rollno = $1, password = $2
	`, rollno, encryptedPassword, chatID)
	if err != nil {
		return err, 20
	}

	return nil, 0
}
