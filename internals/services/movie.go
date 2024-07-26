package services

import (
	"database/sql"
	"errors"
	"fmt"
	constants "movie-rating-api-go/internals"
	"movie-rating-api-go/internals/database"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx"
)

type rating struct {
	Rating float32 `json:"rating"`
}

type MovieCreationReqBody struct {
	Name           string           `json:"name"`
	Description    string           `json:"description"`
	Images         []sql.NullString `json:"images"`
	Release_date   sql.NullTime     `json:"release_date"`
	Director       sql.NullString   `json:"director"`
	Rating_average sql.NullFloat64  `json:"rating_average"`
	Duration       sql.NullInt32    `json:"duration"`
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

// durationFmt: "03:30" 'hh:mm' hh: hour, mm: minute
func convertDurationFormatToMinutes(durationFmt *string) (int, error) {
	// hourMinute: [hh, mm]
	hourMinute := strings.Split(*durationFmt, ":")

	total := 0

	if hrs, err := strconv.Atoi(hourMinute[0]); err != nil {
		return 0, err
	} else {
		total += hrs * 24
	}
	if mins, err := strconv.Atoi(hourMinute[1]); err != nil {
		return 0, err
	} else {
		total += mins
	}

	return total, nil
}

func GetMoviesWithFilter(filters *map[string]interface{}) (*pgx.Rows, error) {
	conn := database.GetConn()

	query := "SELECT * FROM MOVIES WHERE"
	numbering := 1

	var readyFilters []string

	if title, exists := (*filters)["title"]; exists {
		query += fmt.Sprintf(" name = $%d", numbering)

		// make sure (assert) that "title" is a string:
		name, ok := title.(string)
		if !ok {
			return nil, errors.New(constants.ErrSearchFilterTitleIsNotString)
		}
		readyFilters = append(readyFilters, fmt.Sprintf("%%%s%%", name))
		numbering++
	}

	if releaseYear, exists := (*filters)["release_year"]; exists {
		if numbering > 1 {
			query += fmt.Sprintf(" release_year = $%d", numbering)
		} else {
			query += fmt.Sprintf(" AND release_year = $%d", numbering)
		}

		// make sure (assert) that "release_year" is a time.Time:
		releaseDate, ok := releaseYear.(time.Time)
		if !ok {
			return nil, errors.New(constants.ErrSearchFilterReleaseYearIsNotTime)
		}

		readyFilters = append(readyFilters, fmt.Sprintf("%d", releaseDate.Year()))
		numbering++
	}

	if rating, exists := (*filters)["rating"]; exists {
		if numbering > 1 {
			query += fmt.Sprintf(" rating_average > $%d", numbering)
		} else {
			query += fmt.Sprintf(" AND rating_average > $%d", numbering)
		}

		// make sure (assert) that "rating" is a float64:
		ratingAverage, ok := rating.(float64)
		if !ok {
			return nil, errors.New(constants.ErrSearchFilterRatingIsNotFloat64)
		}

		readyFilters = append(readyFilters, fmt.Sprintf("%f", ratingAverage))
		numbering++
	}

	if director, exists := (*filters)["director"]; exists {
		if numbering > 1 {
			query += fmt.Sprintf(" director = $%d", numbering)
		} else {
			query += fmt.Sprintf(" AND director = $%d", numbering)
		}

		// make sure (assert) that "rating" is a float64:
		directorName, ok := director.(string)
		if !ok {
			return nil, errors.New(constants.ErrSearchFilterDirectorIsNotString)
		}

		readyFilters = append(readyFilters, directorName)
		numbering++
	}

	if duration, exists := (*filters)["duration"]; exists {
		if numbering > 1 {
			query += fmt.Sprintf(" duration = $%d", numbering)
		} else {
			query += fmt.Sprintf(" AND duration = $%d", numbering)
		}

		// make sure (assert) that "rating" is a float64:
		durationFormat, ok := duration.(string)
		if !ok {
			return nil, errors.New(constants.ErrSearchFilterDirectorIsNotString)
		}

		durationTime, err := convertDurationFormatToMinutes(&durationFormat)

		if err != nil {
			return nil, err
		}

		readyFilters = append(readyFilters, fmt.Sprintf("%d", durationTime))
		numbering++
	}

	query += ";"

	// for int i = 0; i < numbering
	rows, err := conn.Query(query, readyFilters)

	if err != nil {
		return nil, err
	}

	return rows, nil
}


/*
movie table definition:

CREATE TABLE IF NOT EXISTS MOVIES (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description TEXT NOT NULL,
    images TEXT[],
    release_date DATE,
    director VARCHAR(100),
    rating_average NUMERIC(3, 2) DEFAULT 0.0,
    duration INT -- duration in minutes
);
*/

func AddMovie(movie *MovieCreationReqBody) {
	query := "INSERT INTO MOVIES VALUES($1, $2, $3, $4, $5, $6, $7);"

	conn := database.GetConn()

	if tag, err := conn.Exec(query, ); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Error while returning connection in movie service: %v", err)
		return
	}
} 