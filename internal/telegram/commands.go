package telegram

import (
	"database/sql"
	"log"
	"psghostelbot/internal/db"
	"strings"
)

const (
	registration_tip = "/register <rollno> <password>\n/register 22AB123 myPassword"
)

func handleBotCommand(psg *TelegramBot, msg *Message, db *sql.DB) {
	args := strings.Fields(strings.TrimSpace(msg.Text))
	if len(args) == 0 {
		return
	}
	command := args[0]

	var reply string
	var err error

	switch command {
	case "/start":
		reply, err = handleCommandStart()
	case "/register":
		reply, err = handleCommandRegister(msg, db, args, psg)
	case "/unregister":
		reply, err = handleCommandUnregister(msg, db)
	case "/help":
		reply, err = handleCommandHelp()
	default:
		reply, err = handleUnknownCommand()
	}

	if err != nil {
		log.Printf("Error handling command %s: %v", command, err)
	} else {
		err = psg.sendMessage(reply, &msg.Chat.ID)
		if err != nil {
			log.Printf("Error sending reply for command %s: %v", command, err)
		} else {
			if command == "/register" && len(args) == 3 {
				log.Printf("User %d registered with roll number %s", msg.Chat.ID, strings.Fields(msg.Text)[1])
				psg.deleteMessage(msg.Chat.ID, msg.MessageID)
			}
		}
	}
}

func handleCommandStart() (string, error) {
	message := `
		Welcome to PSG Tech Hostel Bot!
		Register yourself as follows:

	` + registration_tip

	return message, nil
}

func handleCommandRegister(msg *Message, db *sql.DB, args []string, psg *TelegramBot) (string, error) {
	var reply string
	var err error

	if len(args) == 3 {
		rollno := args[1]
		password := args[2]

		reply, err = getRegisterReplyAndError(db, rollno, password, msg.Chat.ID, psg)
		if err != nil {
			log.Print(err)
		}
	} else {
		reply = "Invalid format! Try again as follows:\n" + registration_tip
	}

	return reply, err
}

func handleCommandUnregister(msg *Message, database *sql.DB) (string, error) {
	err := db.DeleteStudent(database, msg.Chat.ID)
	if err != nil {
		log.Printf("Error deleting chat_id %d from database: %v", msg.Chat.ID, err)
		return "An error occurred while processing your request.", err
	} else {
		log.Printf("User %d unregistered successfully.", msg.Chat.ID)
	}

	reply := "You have been unregistered successfully."
	return reply, nil
}

func handleCommandHelp() (string, error) {
	message := `
		Available commands:
		/start - Start the bot and get registration instructions.
		/register <rollno> <password> - Register yourself to receive booked QR codes.
		/unregister - Unregister yourself from receiving booked QR codes.
		/help - Show this help message.
	`

	return message, nil
}

func handleUnknownCommand() (string, error) {
	message := `
		Unknown command. Click /help to open list of available commands.
	`

	return message, nil
}
