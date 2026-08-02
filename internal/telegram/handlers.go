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
	existingChat := db.QueryRow("SELECT chat_id FROM students WHERE rollno = $1", rollno)
	var existingChatID int64
	err := existingChat.Scan(&existingChatID)
	if err != nil && err != sql.ErrNoRows {
		return err, 20
	}
	if err == nil && existingChatID != chatID {
		return fmt.Errorf("rollno %s is already registered with another chat ID", rollno), 10
	}

	_, err = db.Exec(`
		INSERT INTO students (rollno, password, chat_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (rollno) DO UPDATE SET password = $2, chat_id = $3
	`, rollno, encryptedPassword, chatID)
	if err != nil {
		return err, 20
	}

	return nil, 0
}
