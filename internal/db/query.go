package db

import (
	"database/sql"
	"fmt"
)

type Student struct {
	Rollno      string
	EncPassword string
	ChatID      int64
}

func GetAllUsers(database *sql.DB) ([]Student, error) {
	rows, err := database.Query("SELECT rollno, password, chat_id FROM student")
	if err != nil {
		return nil, fmt.Errorf("UNABLE TO QUERY USERS : %w", err)
	}
	defer rows.Close()

	var students []Student
	for rows.Next() {
		var s Student
		if err := rows.Scan(&s.Rollno, &s.EncPassword, &s.ChatID); err != nil {
			continue // Skip this row if there's an error scanning
		}
		students = append(students, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("UNABLE TO ITERATE USERS : %w", err)
	}

	return students, nil
}

func DeleteStudent(database *sql.DB, chatID int64) error {
	result, err := database.Exec("DELETE FROM student WHERE chat_id=$1", chatID)
	if err != nil {
		return fmt.Errorf("UNABLE TO DELETE USER : %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("UNABLE TO GET ROWS AFFECTED : %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("NO USER FOUND WITH CHAT ID: %d", chatID)
	}

	return nil
}
