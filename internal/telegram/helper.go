package telegram

import (
	"database/sql"
	"fmt"
	pq "github.com/lib/pq"
	"os"
	"psghostelbot/internal/crypt"
	"psghostelbot/internal/scraper"
	"strings"
)

func dedent(s string) string {
	s = strings.TrimPrefix(s, "\n")
	lines := strings.Split(s, "\n")
	min := -1
	for _, line := range lines {
		trimmed := strings.TrimLeft(line, " \t")
		if trimmed == "" {
			continue
		}
		indent := len(line) - len(trimmed)
		if min == -1 || indent < min {
			min = indent
		}
	}
	if min > 0 {
		for i, line := range lines {
			if len(line) >= min {
				lines[i] = line[min:]
			}
		}
	}
	return strings.Join(lines, "\n")
}

func insertStudent(db *sql.DB, rollno, password string, chatID int64) (string, error) {
	var botReply string

	_, err := db.Exec(`
		INSERT INTO student (rollno, password, chat_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (chat_id) DO UPDATE SET rollno=$1, password = $2
	`, rollno, password, chatID)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			botReply = fmt.Sprintf("%s is already registered by another user!", rollno)
			return botReply, fmt.Errorf("%s is already registered with another chat ID", rollno)
		} else {
			botReply = "An error occurred while registering! Please try again."
			return botReply, fmt.Errorf("Database error: %v", err)
		}
	}

	return botReply, nil
}

func getRegisterReplyAndError(db *sql.DB, rollno, password string, chatID int64, psg *TelegramBot) (string, error) {
	var reply string = "Invalid format! Try again as follows:\n" + registration_tip

	username, err := scraper.GetUserIfExists(rollno, password)
	if err != nil {
		reply = "An error occurred while registering! Please try again."
		return reply, err
	} else if username != "" {
		encryptedPassword, err := crypt.EncryptPassword(password, os.Getenv("AES_KEY"))
		if err != nil {
			reply = "An error occurred while encrypting your data. Please try again."
			return reply, err
		} else {
			replyMsg, err := insertStudent(db, rollno, encryptedPassword, chatID)
			if err != nil {
				reply = replyMsg
				return reply, err
			} else {
				reply = fmt.Sprintf("Hello %s! You have been successfully registered to receive your booked QR codes!\n", username)
				RunScraperForUser(rollno, password, chatID, psg)
			}
		}
	} else {
		reply = "Invalid rollno or password! Please try again."
	}

	return reply, nil
}
