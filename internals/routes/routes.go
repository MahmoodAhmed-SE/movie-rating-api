package routes

import (
	"net/http"

	"movie-rating-api-go/internals/handlers"
	"movie-rating-api-go/internals/middlewares"
)

func SetupRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/register-user", handlers.UserRegistration)
	mux.HandleFunc("/api/v1/login-user", handlers.UserLogin)

	mux.Handle("/api/v1/retrieve-movies", middlewares.JWTAuthorization(http.HandlerFunc(handlers.MoviesRetrievel)))
	return mux
}
