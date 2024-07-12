package initialdata

import (
	"encoding/json"
	"fmt"
	"log"
	constants "movie-rating-api-go/internals"
	"movie-rating-api-go/internals/database"
	"movie-rating-api-go/internals/models"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx"
	"github.com/joho/godotenv"
)

type ExtraDataResponseBody struct {
	ReleaseDate string `json:"release_date"`
	Director    string `json:"director"`
}

func StartMovieStartupInsertions() {
	var movies []models.Movie

	moviesBytes, err := os.ReadFile("internals/initial_data/movies.json")
	if err != nil {
		log.Printf("Error reading movies file: %s", err)
	}

	if parsingErr := json.Unmarshal(moviesBytes, &movies); parsingErr != nil {
		log.Printf("Error parsing movies json into movies array of Movie: %v", parsingErr)
	}

	if envErr := godotenv.Load(); envErr != nil {
		log.Printf("Error loading env: %v", envErr)
	}
	var conn *pgx.Conn = database.GetConn()

	if assurance, selectErr := conn.Query("SELECT * FROM MOVIES;"); selectErr != nil {
		log.Fatalf("Error in select movies query: %v", selectErr)
		return
	} else {
		defer assurance.Close()

		number_of_rows := 0
		for assurance.Next() {
			number_of_rows++
			if number_of_rows >= 50 {
				break
			}
		}

		if number_of_rows == 50 {
			log.Println("Initial rows are already added. There is no need for movies startup insertions.")
			return
		}
	}

	for _, movie := range movies {
		resp, fetchErr := http.Get(fmt.Sprintf("%s%d", os.Getenv(constants.EnvMovieAPI), movie.Id))
		if fetchErr != nil {
			log.Printf("Error fetching movies info from an external movie api: %v", fetchErr)

		}
		defer resp.Body.Close()

		var respBody ExtraDataResponseBody

		decoder := json.NewDecoder(resp.Body)
		if decodeErr := decoder.Decode(&respBody); decodeErr != nil {
			log.Printf("Error decoding movies extra data: %v", decodeErr)
		}

		movie.Director.String = respBody.Director
		movie.Director.Valid = true

		t, timeErr := time.Parse("02-01-2006", respBody.ReleaseDate)
		if timeErr != nil {
			log.Printf("error time %v", timeErr)
		}

		movie.Release_date.Time = t
		movie.Director.Valid = true
		/*
			name
			description
			images
			release_date
			director
			rating_average
			duration
		*/

		_, execErr := conn.Exec("INSERT INTO MOVIES VALUES (DEFAULT, $1, $2, NULL, $3, $4);", movie.Name, movie.Description, movie.Release_date, movie.Director)
		if execErr != nil {
			log.Fatalf("Error executing movie insertion: %v", execErr)
		}
	}

}
