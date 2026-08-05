package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

func NewConnection(connString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("UNABLE TO OPEN DATABASE: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("UNABLE TO PING DATABASE: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(15 * time.Minute)

	return db, nil
}
