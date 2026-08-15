package telegram

import "fmt"

type AuthenticationError struct {
	Message string
	Err     error
}

type DatabaseError struct {
	Query string
	Err   error
}

func (e *AuthenticationError) Unwrap() error {
	return e.Err
}

func (e *AuthenticationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Message, e.Err)
}

func (e *DatabaseError) Unwrap() error {
	return e.Err
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("%s: %s", e.Query, e.Err)
}
