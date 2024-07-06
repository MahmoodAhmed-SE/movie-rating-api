package routes

import (
	"net/http"

	"movie-rating-api-go/internals/handlers"
	"movie-rating-api-go/internals/middlewares"
)

func SetupRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/register-user", handlers.HandleUserRegistration)
	mux.Handle("/api/retrieve-movies", middlewares.JWTAuthorization(http.HandlerFunc(handlers.HandleMoviesRetrievel)))

	return mux
}
