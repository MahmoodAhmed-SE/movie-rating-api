package movie

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	constants "movie-rating-api-go/internals"
	"movie-rating-api-go/internals/services"
	"net/http"

	"github.com/jackc/pgx"
)

type RatingRequestBody struct {
	MovieId int            `json:"movie_id"`
	Rating  float32        `json:"rating"`
	Review  sql.NullString `json:"review"`
}

func RatingMovie(w http.ResponseWriter, r *http.Request) {
	user_id, ok := r.Context().Value(constants.UserIdKey).(int)
	if !ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println("Error while parsing id passed from authorization middleware to string.")
		return
	}
	decoder := json.NewDecoder(r.Body)
	var data RatingRequestBody
	if err := decoder.Decode(&data); err != nil {
		http.Error(w, "Bad request", http.StatusInternalServerError)
		log.Printf("Error while decoding req body into data (RatingRequestBody): %v.", err)
		return
	}

	rating, err := services.GetUserRating(user_id, data.MovieId)

	// Check whether user has already rated this movie or an error occured while querying
	if err == nil {
		http.Error(w, fmt.Sprintf("Movie has already been rated with %v/10", rating), http.StatusConflict)
		return
	} else if err != pgx.ErrNoRows {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Error while querying user rating: %v.", err)
		return
	}

	rateErr := services.RateMovie(user_id, data.MovieId, data.Rating, data.Review)

	if rateErr != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Error while querying inserting a rating: %v.", err)
		return
	}
}
