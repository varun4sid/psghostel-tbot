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

const (
	registration_tip = "/register <rollno> <password>\n/register 22AB123 myPassword"
)

func handleBotCommand(psg *TelegramBot, msg *Message, db *sql.DB) {
	args := strings.Split(strings.TrimSpace(msg.Text), " ")
	command := args[0]

	var err error
	switch command {
	case "/start":
		err = handleCommandStart(psg, &msg.Chat.ID)
	case "/register":
		err = handleCommandRegister(psg, msg, db)
	default:
		err = handleUnknownCommand(psg, msg)
	}

	if err != nil {
		slog.Error("", err)
	}
}

func handleCommandStart(psg *TelegramBot, chatID *int64) error {
	message := `
		Welcome to PSG Tech Hostel Bot!
		Register yourself as follows:

	` + registration_tip

	return psg.sendMessage(message, chatID)
}

func handleCommandRegister(psg *TelegramBot, msg *Message, db *sql.DB) error {
	args := strings.Fields(strings.TrimSpace(msg.Text))
	var opError error

	reply := "Invalid format! Try again as follows:\n" + registration_tip
	if len(args) == 3 {
		rollno := args[1]
		password := strings.ToUpper(args[2])

		username, err := scraper.GetUserIfExists(rollno, password)
		if err != nil {
			reply = "An error occurred while registering! Please try again."
			opError = err
		} else if username != "" {
			encryptedPassword, err := crypt.EncryptPassword(password, os.Getenv("AES_KEY"))
			if err != nil {
				reply = "An error occurred while encrypting your data. Please try again."
				opError = err
			} else {
				replyMsg, err := insertStudent(db, rollno, encryptedPassword, msg.Chat.ID)
				if err != nil {
					reply = replyMsg
					opError = err
				} else {
					reply = fmt.Sprintf("Hello %s! You have been successfully registered!", username)
				}
			}
		} else {
			reply = "Invalid rollno or password! Please try again."
		}
	}

	if sendError := psg.sendMessage(reply, &msg.Chat.ID); sendError != nil {
		return sendError
	}

	if deleteError := psg.deleteMessage(msg.Chat.ID, msg.MessageID); deleteError != nil {
		return deleteError
	}

	return opError
}

func handleUnknownCommand(psg *TelegramBot, msg *Message) error {
	message := `
		Unknown command. Click /help to open list of available commands.
	`

	return psg.sendMessage(message, &msg.Chat.ID)
}
