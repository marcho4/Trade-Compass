package infrastructure

import (
	"errors"
	"fmt"
)

type DbError struct {
	Message      string
	RowsAffected int
}

func (d *DbError) Error() string {
	return fmt.Sprintf("Database error: %s, rows affected: %d", d.Message, d.RowsAffected)
}

func (d *DbError) Unwrap() error {
	return errors.New(d.Message)
}

func NewDbError(message string, rowsAffected int) *DbError {
	return &DbError{
		Message:      message,
		RowsAffected: rowsAffected,
	}
}
