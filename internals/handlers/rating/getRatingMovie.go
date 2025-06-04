package rating

import (
	"encoding/json"
	"log"
	"movie-rating-api-go/internals/database"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
)

type RatingResponse struct {
	Rating float32 `json:"rating"`
}

func GetRatingMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	movieId := vars["movieId"]

	conn := database.GetConn()

	row := conn.QueryRow("SELECT rating FROM RATINGS WHERE movie_id = $1;", movieId)

	var ratingResp RatingResponse

	if err := row.Scan(&ratingResp.Rating); err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Rating not found", http.StatusNotFound)
			return
		}
		log.Printf("Error while scanning rating row to rating response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(ratingResp); err != nil {
		log.Printf("Error while decoding rating response to json: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
