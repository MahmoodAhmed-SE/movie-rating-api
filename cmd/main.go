package main

import (
	"fmt"
	"log"
	"movie-rating-api-go/internals/config"
	"movie-rating-api-go/internals/database"
	initialdata "movie-rating-api-go/internals/initial_data"
	"movie-rating-api-go/internals/routes"
	"net/http"
)

func main() {
	config := config.LoadConfig()

	if err := database.ConnectPostgres(config.PgxConfig); err != nil {
		log.Fatalf("Error connecting to moviesdb: %v", err)
	}

	conn := database.GetConn()
	defer conn.Close()

	//
	initialdata.StartMovieStartupInsertions()

	mux := routes.SetupRoutes()

	log.Printf("Server is starting up on port: %d", config.ServerConfig.Port)

	address := fmt.Sprintf("%s:%d", config.ServerConfig.IP, config.ServerConfig.Port)
	if err := http.ListenAndServe(address, mux); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}

}
