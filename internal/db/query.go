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
