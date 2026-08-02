package telegram

import (
	"fmt"
	"log/slog"
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
	fmt.Printf("chatID : %v", chatID)
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

		botReply = "Roll : " + rollno
		slog.Info("/register", "rollno", rollno)
	}

	if err := psg.sendJSON("sendMessage", SendMessagePayload{
		ChatID:    msg.Chat.ID,
		Text:      botReply,
		ParseMode: "Markdown",
	}); err != nil {
		slog.Error("Failed to send message", "error", err)
	}
}
