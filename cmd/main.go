package main

import (
	"fmt"
	"log"
	"movie-rating-api-go/internals/config"
	"movie-rating-api-go/internals/database"
	"movie-rating-api-go/internals/routes"
	"net/http"
)

type USER struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	config := config.LoadConfig()

	if err := database.ConnectPostgres(config.PgxConfig); err != nil {
		log.Fatalf("Error connecting to moviesdb: %v", err)
	}

	conn := database.GetConn()
	defer conn.Close()

	routes.SetupRoutes()

	log.Printf("Server is starting up on port: %d", config.ServerConfig.Port)

	address := fmt.Sprintf("%s:%d", config.ServerConfig.IP, config.ServerConfig.Port)
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}

}
