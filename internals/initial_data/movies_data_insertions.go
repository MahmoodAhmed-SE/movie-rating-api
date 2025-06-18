package initialdata

import (
	"encoding/json"
	"log"
	"movie-rating-api-go/internals/database"
	"movie-rating-api-go/internals/models"
	"movie-rating-api-go/internals/services"
	"os"

	"github.com/joho/godotenv"
)

func StartMovieStartupInsertions() {
	if envErr := godotenv.Load(); envErr != nil {
		log.Printf("Error loading env: %v", envErr)
	}

	// simple logic to check if there is already enough sample rows..
	conn := database.GetConn()

	if rows, selectErr := conn.Query("SELECT * FROM MOVIES;"); selectErr != nil {
		log.Fatalf("Error in select movies query: %v", selectErr)
		return
	} else {
		defer rows.Close()

		number_of_rows := 0
		for rows.Next() {
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
		log.Fatalf("Error reading movies file: %s", err)
	}

	if err = json.Unmarshal(moviesBytes, &movies); err != nil {
		log.Fatalf("Error parsing movies json into movies array of Movie: %v", err)
	}

	log.Println("Initializing movies for MOVIES table")
	// insert rows to movies
	for _, movie := range movies {
		descriptionVec, err := services.GetGrpcEmbeddingResp(movie.Description)
		if err != nil {
			log.Fatalf("Error executing movie insertion grpc: %v", err)
		}

		_, err = conn.Exec("INSERT INTO MOVIES VALUES (DEFAULT, $1, $2, $3);", movie.Name, movie.Description, descriptionVec)
		if err != nil {
			log.Fatalf("Error executing movie insertion: %v", err)
		}
	}

	log.Println("Finished Initializing 50 movies for MOVIES table!")
}
