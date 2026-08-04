package telegram

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"psghostelbot/internal/crypt"
	"psghostelbot/internal/scraper"
	"strings"
)

func handleBotCommand(psg *TelegramBot, msg *Message, db *sql.DB) {
	args := strings.Split(strings.TrimSpace(msg.Text), " ")
	command := args[0]

	switch command {
	case "/start":
		handleCommandStart(psg, &msg.Chat.ID)
	case "/register":
		handleCommandRegister(psg, msg, db)
	default:
		handleUnknownCommand(psg, msg)
	}
}

func handleCommandStart(psg *TelegramBot, chatID *int64) {
	message := `
		Welcome to PSG Tech Hostel Bot!
		Register yourself as follows:

		/register <rollno> <password>
		/register 22AB123 myPassword
	`

	if err := psg.sendMessage(message, chatID); err != nil {
		slog.Error("Failed to send message", "error", err)
	}
}

func handleCommandRegister(psg *TelegramBot, msg *Message, db *sql.DB) {
	args := strings.Split(strings.TrimSpace(msg.Text), " ")

	var botReply string

	if len(args) == 3 {
		rollno := args[1]
		password := strings.ToUpper(args[2])

		username, err := scraper.GetUserIfExists(rollno, password)
		if err != nil {
			fmt.Printf("\nError validating user credentials : %v\n", err)
		}

		if username != "" {
			encryptedPassword, err := crypt.EncryptPassword(password, os.Getenv("AES_KEY"))
			if err != nil {
				slog.Error("Failed to encrypt password", "error", err)
				botReply = "An error occurred while encrypting. Please try again."
			}

			err, errMsg := insertStudent(db, rollno, encryptedPassword, msg.Chat.ID)
			if err != nil {
				botReply = errMsg
			} else {
				botReply = fmt.Sprintf("Hello %s! You have been successfully registered!", username)
			}
		} else {
			botReply = "Invalid credentials! Please try again."
		}
	} else {
		botReply = "Invalid format! Try again."
	}

	if err := psg.sendMessage(botReply, &msg.Chat.ID); err != nil {
		slog.Error("Failed to send message", "error", err)
	}

	if err := psg.deleteMessage(msg.Chat.ID, msg.MessageID); err != nil {
		slog.Error("Failed to delete message", "error", err)
	}
}

func handleUnknownCommand(psg *TelegramBot, msg *Message) {
	message := `
		Unknown command. Click /help to open list of available commands.
	`

	psg.sendMessage(message, &msg.Chat.ID)
}
