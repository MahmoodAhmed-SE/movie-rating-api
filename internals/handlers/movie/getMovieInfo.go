package movie

import (
	"encoding/json"
	"log"
	"movie-rating-api-go/internals/database"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
)

func GetMovieInfo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	movieId := vars["movieId"]

	// TO-DO: validate movieId

	reqMethod := r.Method
	if reqMethod == http.MethodGet {
		conn := database.GetConn()

		row := conn.QueryRow("SELECT name, description FROM MOVIES WHERE id = $1;", movieId)

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
	} else {
		log.Printf("Method Not Allowed: %s", reqMethod)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
