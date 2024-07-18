package movie

import (
	"encoding/json"
	"log"
	"movie-rating-api-go/internals/database"
	"movie-rating-api-go/internals/models"
	"net/http"
)

type movieRetrievelRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

func MoviesRetrievel(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var data movieRetrievelRequestBody
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		log.Printf("Error decoding movies retrievel api request body: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	conn := database.GetConn()

	rows, err := conn.Query("SELECT * FROM MOVIES;")
	if err != nil {
		log.Printf("Error retrieving all movies: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
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
			log.Printf("Error scanning to movie model %v", scanningErr)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		movies = append(movies, movie)
	}

	encoder := json.NewEncoder(w)
	if jsonEncodingErr := encoder.Encode(movies); jsonEncodingErr != nil {
		log.Printf("Error encoding movies array/slice to response writer %v", jsonEncodingErr)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
