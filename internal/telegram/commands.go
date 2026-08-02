package telegram

import (
	"fmt"
	"log/slog"
	"psghostelbot/internal/scraper"
	"strings"
)

func handleBotCommand(psg *TelegramBot, msg *Message) {
	args := strings.Split(strings.TrimSpace(msg.Text), " ")
	command := args[0]

	switch command {
	case "/start":
		handleCommandStart(psg, &msg.Chat.ID)
	case "/register":
		handleCommandRegister(psg, msg)
	}
}

func handleCommandStart(psg *TelegramBot, chatID *int64) {
	if err := psg.sendJSON("sendMessage", SendMessagePayload{
		ChatID: *chatID,
		Text: dedent(`
			Welcome to PSG Tech Hostel Bot!
			Register yourself using the following command:
			/register <rollno> <password>
			Eg.
			/register 22AB123 myPassword
		`),
		ParseMode: "Markdown",
	}); err != nil {
		slog.Error("Failed to send message", "error", err)
	}
}

func handleCommandRegister(psg *TelegramBot, msg *Message) {
	args := strings.Split(strings.TrimSpace(msg.Text), " ")

	var botReply string

	if len(args) != 3 {
		botReply = dedent(`
			Invalid format! Try again!
		`)
	} else {
		rollno := args[1]
		password := strings.ToUpper(args[2])

		isValid, username, err := scraper.CheckValidCredentials(rollno, password)
		if err != nil {
			fmt.Printf("\nError validating user credentials : %v\n", err)
		}
		if isValid {
			botReply = fmt.Sprintf("Hello %s! You have been successfully registered!", username)
		} else {
			botReply = "Invalid credentials. Please try again."
		}
	}

	if err := psg.sendJSON("sendMessage", SendMessagePayload{
		ChatID:    msg.Chat.ID,
		Text:      botReply,
		ParseMode: "Markdown",
	}); err != nil {
		slog.Error("Failed to send message", "error", err)
	}
}
