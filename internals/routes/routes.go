package routes

import (
	"net/http"

	"movie-rating-api-go/internals/handlers/auth"
	"movie-rating-api-go/internals/handlers/movie"
	"movie-rating-api-go/internals/handlers/rating"
	"movie-rating-api-go/internals/handlers/streaming"
	"movie-rating-api-go/internals/handlers/watchlist"
	"movie-rating-api-go/internals/middlewares"

	"github.com/gorilla/mux"
)

func SetupRoutes() http.Handler {
	jwtAuth := middlewares.JWTAuthorization
	router := mux.NewRouter()

	// Authentication Routes
	router.HandleFunc("/api/v1/register-user", auth.UserRegistration).Methods("POST")
	router.HandleFunc("/api/v1/login-user", auth.UserLogin).Methods("POST")

	// Movie Retrieval & Info
	router.Handle("/api/v1/movies", jwtAuth(http.HandlerFunc(movie.MoviesRetrievel))).Methods("GET")
	router.Handle("/api/v1/movies/{movieId:[0-9]+}", jwtAuth(http.HandlerFunc(movie.GetMovieInfo))).Methods("GET")
	router.Handle("/api/v1/search", jwtAuth(http.HandlerFunc(movie.SearchMovie))).Methods("GET")

	// Movie Rating
	router.Handle("/api/v1/movies-rating/{movieId:[0-9]+}", jwtAuth(http.HandlerFunc(rating.GetRatingMovie))).Methods("GET")
	router.Handle("/api/v1/movies-rating", jwtAuth(http.HandlerFunc(rating.RatingMovie))).Methods("POST")

	// Movie Chat
	router.Handle("/api/v1/chat-on-movie", jwtAuth(http.HandlerFunc(movie.ChatOnMovie))).Methods("POST")
	router.Handle("/api/v1/chat-on-movie/{movieId:[0-9]+}", jwtAuth(http.HandlerFunc(movie.ChatOnMoviePathQuery))).Methods("GET")

	// Watchlist
	router.Handle("/api/v1/watchlist", jwtAuth(http.HandlerFunc(watchlist.PostWatchlist))).Methods("POST")
	router.Handle("/api/v1/watchlist", jwtAuth(http.HandlerFunc(watchlist.GetWatchlist))).Methods("GET")

	// Real-Time Features (WebSocket/RPC)
	router.Handle("/api/v1/ws" /*jwtAuth(*/, http.HandlerFunc(streaming.WSConnHandler)). /*)*/ Methods("GET")

	return router
}
