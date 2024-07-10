package services

import (
	"movie-rating-api-go/internals/database"
	"movie-rating-api-go/internals/models"
)

func GetUser(username string) (models.USER, error) {
	var user models.USER

	conn := database.GetConn()
	row := conn.QueryRow("SELECT * FROM USERS WHERE username = $1;", username)

	if err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return user, err
	}

	return user, nil
}
