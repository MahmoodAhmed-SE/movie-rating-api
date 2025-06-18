package movie

import (
	"encoding/json"
	"log"
	"movie-rating-api-go/internals/database"
	"movie-rating-api-go/internals/models"
	"movie-rating-api-go/internals/services"
	"net/http"
	"strings"
)

type SearchParams struct {
	MovieName  string `json:"movie_name"`
	PageNumber int    `json:"page_number"`
}

func SearchMovie(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var data SearchParams

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Printf("Error while decoding request body data to SearchParams struct: %v", err)
		return
	}

	name := strings.TrimSpace(data.MovieName)

	// rows limit for search result for each page
	const rowsLimit int = 10

	// if page number is 1 then it should list starting by 0, otherwize, (page number - 1) * rowsLimit
	offset := 0
	if data.PageNumber > 0 {
		data.PageNumber--
		offset = data.PageNumber * rowsLimit
	}

	nameVectorized, err := services.GetGrpcEmbeddingResp(name)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error vectorizing name: %v", err)
		return
	}

	conn := database.GetConn()
	rows, err := conn.Query("SELECT id, name, description FROM movies ORDER BY description_vec <-> $1 LIMIT $2 OFFSET $3;", nameVectorized, rowsLimit, offset)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error doing select query: %v", err)
		return
	}

	defer rows.Close()

	var movies []models.Movie

	for rows.Next() {
		var moviesItem models.Movie

		if err := rows.Scan(&moviesItem.Id, &moviesItem.Name, &moviesItem.Description); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("Error doing scanning returned row to models.Movie: %v", err)
			return
		}

		movies = append(movies, moviesItem)
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(movies); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error while encoding []models.Movie movies to json: %v", err)
		return
	}
}
