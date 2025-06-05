package watchlist

import (
	"encoding/json"
	"github.com/jackc/pgx"
	"log"
	constants "movie-rating-api-go/internals"
	"movie-rating-api-go/internals/database"
	"movie-rating-api-go/internals/models"
	"net/http"
)

func GetWatchlist(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userId := r.Context().Value(constants.UserIdKey)

	if userId == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Printf("user ['user_id'] in request context is not available even though reached this code in [GetWatchlist] func: %v", userId)
		return
	}

	conn := database.GetConn()

	rows, err := conn.Query("SELECT movie_id, user_id, created_at FROM WATCHLIST WHERE user_id = $1;", userId)

	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Not Found", http.StatusNotFound)
			log.Printf("info: No watchlist record found with user_id: %v", userId)
			return
		}

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("An Error occurred while querying for watchlist with user_id %v: %v", userId, err)
		return
	}

	var watchlist []models.WatchlistResponse

	for rows.Next() {
		var watchlistItem models.WatchlistResponse

		if err := rows.Scan(&watchlistItem.UserId, &watchlistItem.MovieId, &watchlistItem.CreatedAt); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("An Error occurred while scanning the watch list rows to [models.WatchlistResponse] instance: %v", err)
			return
		}

		watchlist = append(watchlist, watchlistItem)
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(watchlist); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("An Error occurred while encoding the watchlist to json: %v", err)
		return
	}
}
