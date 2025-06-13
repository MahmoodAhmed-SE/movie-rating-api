package main

import (
	"os"
	"fmt"
	"log"
	"movie-rating-api-go/internals/config"
	"movie-rating-api-go/internals/database"
	initialdata "movie-rating-api-go/internals/initial_data"
	"time"
	"movie-rating-api-go/internals/routes"
	"net/http"
	"github.com/gorilla/handlers"
)

func main() {
	config := config.LoadConfig()

	if err := database.ConnectPostgres(config.PgxConfig); err != nil {
		log.Fatalf("Error connecting to moviesdb: %v", err)
	}

	conn := database.GetConn()
	defer conn.Close()

	initialdata.StartMovieStartupInsertions()

	muxRouter := routes.SetupRoutes()

	// adding logging middleware
	loggedRouter := handlers.LoggingHandler(os.Stdout, muxRouter) 

	log.Printf("Server is starting up on port: %d", config.ServerConfig.Port)

	address := fmt.Sprintf("%s:%d", config.ServerConfig.IP, config.ServerConfig.Port)

	s := &http.Server{
		Addr: address,
		Handler: loggedRouter,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("Error: Server shutdown or unable to start: %s", err)
	}

}
