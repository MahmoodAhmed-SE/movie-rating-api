package services

import (
	"database/sql"
	"movie-rating-api-go/internals/database"
)

type rating struct {
	Rating float32 `json:"rating"`
}

func GetUserRating(userId int, movieId int) (float32, error) {
	conn := database.GetConn()

	row := conn.QueryRow(`
		SELECT rating FROM RATINGS 
		WHERE user_id = $1 AND movie_id = $2;`,
		userId, movieId)

	var r rating
	if err := row.Scan(&r.Rating); err != nil {
		return 0, err
	}

	return r.Rating, nil
}

func RateMovie(userId int, movieId int, rate float32, review sql.NullString) error {
	conn := database.GetConn()

	_, err := conn.Exec(`
		INSERT INTO RATINGS
		VALUES(DEFAULT, $1, $2, $3, $4, DEFAULT);`,
		userId, movieId, rate, review)

	if err != nil {
		return err
	}

	return nil
}
