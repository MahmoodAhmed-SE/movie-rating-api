package initialdata

import (
	"encoding/json"
	"log"
	"movie-rating-api-go/internals/database"
	"movie-rating-api-go/internals/models"
	"os"

	"github.com/jackc/pgx"
	"github.com/joho/godotenv"
)

func StartMovieStartupInsertions() {
	if envErr := godotenv.Load(); envErr != nil {
		log.Printf("Error loading env: %v", envErr)
	}

	// simple logic to check if there is already enough sample rows..
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

	// get movies data from json file and assign them to movies with type [[]models.Movie]..
	var movies []models.Movie

	moviesBytes, err := os.ReadFile("internals/initial_data/movies.json")
	if err != nil {
		log.Printf("Error reading movies file: %s", err)
	}

	if parsingErr := json.Unmarshal(moviesBytes, &movies); parsingErr != nil {
		log.Printf("Error parsing movies json into movies array of Movie: %v", parsingErr)
	}

	log.Println("Initializing 50 movies for MOVIES table")
	// insert rows to movies
	for _, movie := range movies {
		/*
			name
			description
			images
		*/

		_, execErr := conn.Exec("INSERT INTO MOVIES VALUES (DEFAULT, $1, $2);", movie.Name, movie.Description)
		if execErr != nil {
			log.Fatalf("Error executing movie insertion: %v", execErr)
			break
		}
	}
	log.Println("Finished Initializing 50 movies for MOVIES table!")

}
