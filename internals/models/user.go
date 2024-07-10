package models

import (
	"database/sql"
	"time"
)

type USER struct {
	ID        int            `json:"id"`
	Username  string         `json:"username"`
	Password  string         `json:"password"`
	Email     sql.NullString `json:"email"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
