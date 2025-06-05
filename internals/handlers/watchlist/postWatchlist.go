package watchlist

import (
	"encoding/json"
	"log"
	constants "movie-rating-api-go/internals"
	"movie-rating-api-go/internals/database"
	"net/http"
)

type ReqBody struct {
	MovieId int `json:"movie_id"`
}

func PostWatchlist(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userId := r.Context().Value(constants.UserIdKey)

	if userId == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Printf("user ['user_id'] in request context is not available even though reached this code in [GetWatchlist] func: %v", userId)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var reqBody ReqBody

	if err := decoder.Decode(&reqBody); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Printf("Error while decoding req body of postwatchlist func: %v", err)
		return
	}

	movieId := reqBody.MovieId

	conn := database.GetConn()

	_, err := conn.Exec("INSERT INTO WATCHLIST(id, movie_id, user_id, created_at) VALUES(DEFAULT, $1, $2, DEFAULT);", movieId, userId)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("An Error occurred while inserting a watchlist : %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Watchlist item added successfully!"))
}
