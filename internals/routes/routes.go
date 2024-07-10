package routes

import (
	"net/http"

	"movie-rating-api-go/internals/handlers"
	"movie-rating-api-go/internals/middlewares"
)

func SetupRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/register-user", handlers.UserRegistration)
	mux.Handle("/api/retrieve-movies", middlewares.JWTAuthorization(http.HandlerFunc(handlers.MoviesRetrievel)))

	return mux
}
