package models

import (
	"database/sql"
	"time"
)

type Rating struct {
	Id         int            `json:"id"`         // Not null
	User_id    sql.NullString `json:"user_id"`    // Can be null
	Movie_id   int            `json:"movie_id"`   // Not null
	Rating     float32        `json:"rating"`     // Not null
	Created_at time.Time      `json:"created_at"` // Default to current timestamp
}
