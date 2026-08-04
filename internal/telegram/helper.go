package telegram

import (
	"database/sql"
	"fmt"
	"strings"

	pq "github.com/lib/pq"
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

func insertStudent(db *sql.DB, rollno, password string, chatID int64) (error, string) {
	_, err := db.Exec(`
		INSERT INTO students (rollno, password, chat_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (chat_id) DO UPDATE SET rollno=$1, password = $2
	`, rollno, password, chatID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return fmt.Errorf("Rollno %s is already registered with another chat ID", rollno),
					"Rollno %s is already registered by another user. Please contact support if this is an error."
			}
		}
	}

	return nil, "Student registered successfully."
}
