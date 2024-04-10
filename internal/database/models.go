// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Car struct {
	ID              uuid.UUID
	CreatedAt       time.Time
	UpdatedAt       time.Time
	RegNum          string
	Mark            string
	Model           string
	Year            int32
	OwnerName       string
	OwnerSurname    string
	OwnerPatronymic sql.NullString
}
