package models

import (
	"database/sql"
)

type Movie struct {
	Id             int              `json:"id"`
	Name           string           `json:"name"`
	Description    string           `json:"description"`
	Images         []sql.NullString `json:"images"`
	Release_date   sql.NullString   `json:"release_date"`
	Director       sql.NullString   `json:"director"`
	Rating_average sql.NullFloat64  `json:"rating_average"`
	Duration       sql.NullInt32    `json:"duration"`
}
