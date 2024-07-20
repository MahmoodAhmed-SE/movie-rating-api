package movie

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	constants "movie-rating-api-go/internals"
	"movie-rating-api-go/internals/models"
	"movie-rating-api-go/internals/services"
	"net/http"
)

// Genre is not yet implemented due to the lack of table key genre.
// var genres = map[string]interface{}{
// 	"Action":   "",
// 	"Comedey":  "",
// 	"Drama":    "",
// 	"Fantasy":  "",
// 	"Horror":   "",
// 	"Mystery":  "",
// 	"Romance":  "",
// 	"Sci-Fi":   "",
// 	"Thriller": "",
// }

type Filter struct {
	Title       sql.NullString  `json:"title"`
	Genre       sql.NullString  `json:"genre"`
	ReleaseYear sql.NullTime    `json:"release_year"`
	Rating      sql.NullFloat64 `json:"rating"`
	Director    sql.NullString  `json:"director"`
	Duration    sql.NullString  `json:"duration"`
}

func (filter *Filter) PopulateFilters(filters *map[string]interface{}) error {
	isEmpty := true

	if filter.Title.Valid {
		isEmpty = false
		(*filters)["title"] = filter.Title.String
	}
	// if filter.Genre.Valid {
	// 	isEmpty = false
	// 	(*filters)["genre"] = filter.Genre.String
	// }
	if filter.ReleaseYear.Valid {
		isEmpty = false
		(*filters)["release_year"] = filter.ReleaseYear.Time
	}
	if filter.Rating.Valid {
		if filter.Rating.Float64 < 0 || filter.Rating.Float64 > 10 {
			return errors.New(constants.ErrRatingOutOfBounds)
		}

		isEmpty = false
		(*filters)["rating"] = filter.Rating.Float64
	}
	if filter.Director.Valid {
		isEmpty = false
		(*filters)["director"] = filter.Director.String
	}
	if filter.Duration.Valid {
		if len(filter.Duration.String) != 5 {
			return errors.New(constants.ErrIncorrectDurationFormat)
		}

		isEmpty = false
		(*filters)["duration"] = filter.Duration.String
	}

	if isEmpty {
		return errors.New(constants.ErrEmptyFilters)
	} else {
		return nil
	}
}

func populate_movies(movies *[]models.Movie, filters *map[string]interface{}) error {
	rows, err := services.GetMoviesWithFilter(filters)

	if err != nil {
		return err
	}

	defer rows.Close()

	var movie models.Movie
	for rows.Next() {
		if scnErr := rows.Scan(
			&movie.Id,
			&movie.Name,
			&movie.Description,
			&movie.Images,
			&movie.Release_date,
			&movie.Rating_average,
			&movie.Director,
			&movie.Duration,
		); scnErr != nil {
			log.Printf("Error while scanning rows returned from GetMoviesWithFilter into models.Movie: %v", scnErr)
			return err
		}

		(*movies) = append((*movies), movie)
	}

	return nil
}

func SearchMovie(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var filter Filter
	if err := decoder.Decode(&filter); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		log.Printf("Error while decoding filter from request body: %v", err)
		return
	}

	var filters = make(map[string]interface{})

	if err := filter.PopulateFilters(&filters); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		log.Fatalln("Empty search filter.")
		return
	}

	var result_movies []models.Movie

	if filter.Title.Valid {
		populate_movies(&result_movies, &filters)
	}

}
