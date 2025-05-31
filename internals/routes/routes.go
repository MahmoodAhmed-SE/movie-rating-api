package routes

import (
	"net/http"

	"movie-rating-api-go/internals/handlers/auth"
	"movie-rating-api-go/internals/handlers/movie"
	"movie-rating-api-go/internals/middlewares"

	"github.com/gorilla/mux"
)

func SetupRoutes() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/register-user", auth.UserRegistration)
	r.HandleFunc("/api/v1/login-user", auth.UserLogin)

	r.Handle("/api/v1/retrieve-movies", middlewares.JWTAuthorization(http.HandlerFunc(movie.MoviesRetrievel)))
	r.Handle("/api/v1/rate-movie", middlewares.JWTAuthorization(http.HandlerFunc(movie.RatingMovie)))
	r.Handle("/api/v1/search-movie", middlewares.JWTAuthorization(http.HandlerFunc(movie.SearchMovie)))
	r.Handle("/api/v1/get-movie-info", middlewares.JWTAuthorization(http.HandlerFunc(movie.GetMovieInfo)))
	r.Handle("/api/v1/chat-on-movie", middlewares.JWTAuthorization(http.HandlerFunc(movie.ChatOnMovie)))
	r.Handle("/api/v1/chat-on-movie/{movieId}", middlewares.JWTAuthorization(http.HandlerFunc(movie.ChatOnMoviePathQuery)))
	r.Handle("/api/v1/watchlist", middlewares.JWTAuthorization(http.HandlerFunc(movie.Watchlist)))

	return r
}
