package rating

import (
	"encoding/json"
	"log"
	constants "movie-rating-api-go/internals"
	"movie-rating-api-go/internals/services"
	"net/http"

	"github.com/jackc/pgx"
)

type RatingRequestBody struct {
	MovieId int     `json:"movie_id"`
	Rating  float32 `json:"rating"`
}

func RatingMovie(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userId, ok := r.Context().Value(constants.UserIdKey).(int)
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

	_, err := services.GetUserRating(userId, data.MovieId)

	// Check whether user has already rated this movie or an error occured while querying
	if err == nil {
		http.Error(w, "Movie has already been rated!", http.StatusConflict)
		return
	} else if err != pgx.ErrNoRows {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Error while querying user rating: %v.", err)
		return
	}

	rateErr := services.RateMovie(userId, data.MovieId, data.Rating)

	if rateErr != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Error while querying inserting a rating: %v.", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Rating has been submitted!"))
}
