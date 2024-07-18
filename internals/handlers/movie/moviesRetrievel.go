package movie

import (
	"encoding/json"
	"log"
	"movie-rating-api-go/internals/database"
	"movie-rating-api-go/internals/models"
	"net/http"
)

func MoviesRetrievel(w http.ResponseWriter, r *http.Request) {
	r.Body.Close()

	conn := database.GetConn()

	rows, err := conn.Query("SELECT * FROM MOVIES;")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Error retrieving all movies: %v", err)
		return
	}
	defer rows.Close()

	var movies []models.Movie

	for rows.Next() {
		var movie models.Movie
		if scanningErr := rows.Scan(
			&movie.Id,
			&movie.Name,
			&movie.Description,
			&movie.Images,
			&movie.Release_date,
			&movie.Director,
			&movie.Rating_average,
			&movie.Duration,
		); scanningErr != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Printf("Error scanning to movie model %v", scanningErr)
			return
		}
		movies = append(movies, movie)
	}

	encoder := json.NewEncoder(w)
	if jsonEncodingErr := encoder.Encode(movies); jsonEncodingErr != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Error encoding movies array/slice to response writer %v", jsonEncodingErr)
		return
	}
}
