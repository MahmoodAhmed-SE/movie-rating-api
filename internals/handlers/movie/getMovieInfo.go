package movie

import (
	"encoding/json"
	"log"
	"movie-rating-api-go/internals/database"
	"net/http"

	"github.com/jackc/pgx"
)

type GetMovieInfoReqBody struct {
	Movie_id int `json:"movie_id"`
}

func GetMovieInfo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	var reqBody GetMovieInfoReqBody

	if err := decoder.Decode(&reqBody); err != nil {
		log.Printf("Error while decoding get-movie-info request body int GetMovieInfoReqBody struct: %v", err)
		http.Error(w, "Not expected format for request body!", http.StatusBadRequest)
		return
	}

	conn := database.GetConn()

	row := conn.QueryRow("SELECT name, description FROM MOVIES WHERE id = $1;", reqBody.Movie_id)

	var movie struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := row.Scan(&movie.Name, &movie.Description); err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "No movie with such id!", http.StatusBadRequest)
			return
		}

		log.Printf("Error while encoding pgx row to custom made struct var in get-movie-info request: %v", err)
		http.Error(w, "Internal Server Error!", http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)

	if err := encoder.Encode(movie); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Error encoding movie to response writer %v", err)
		return
	}
}
