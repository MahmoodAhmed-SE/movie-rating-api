package routes

import (
	"net/http"

	"movie-rating-api-go/internals/handlers/auth"
	"movie-rating-api-go/internals/handlers/movie"
	"movie-rating-api-go/internals/middlewares"
)

func SetupRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/register-user", auth.UserRegistration)
	mux.HandleFunc("/api/v1/login-user", auth.UserLogin)

	mux.Handle("/api/v1/retrieve-movies", middlewares.JWTAuthorization(http.HandlerFunc(movie.MoviesRetrievel)))
	mux.Handle("/api/v1/rate-movie", middlewares.JWTAuthorization(http.HandlerFunc(movie.RatingMovie)))
	mux.Handle("/api/v1/search-movie", middlewares.JWTAuthorization(http.HandlerFunc(movie.SearchMovie)))
	mux.Handle("/api/v1/create-movie", middlewares.JWTAuthorization(http.HandlerFunc(movie.CreateMovie)))


	return mux
}
