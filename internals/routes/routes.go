package routes

import (
	"net/http"

	"movie-rating-api-go/internals/handlers"
)

func SetupRoutes() {
	http.HandleFunc("/api/register-user", handlers.HandleUserRegistration)
}
