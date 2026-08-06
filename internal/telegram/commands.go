package telegram

import (
	"database/sql"
	"log"
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
		reply, err = handleCommandRegister(msg, db, args)
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
			if command == "/register" {
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

func handleCommandRegister(msg *Message, db *sql.DB, args []string) (string, error) {
	var reply string
	var err error

	if len(args) == 3 {
		rollno := args[1]
		password := strings.ToUpper(args[2])

		reply, err = getRegisterReplyAndError(db, rollno, password, msg.Chat.ID)
		if err != nil {
			log.Print(err)
		}
	} else {
		reply = "Invalid format! Try again as follows:\n" + registration_tip
	}

	return reply, err
}

func handleUnknownCommand() (string, error) {
	message := `
		Unknown command. Click /help to open list of available commands.
	`

	return message, nil
}
