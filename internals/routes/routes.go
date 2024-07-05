package routes

import (
	"log"
	"net/http"
	"time"

	"movie-rating-api-go/internals/handlers"
)

type Logger struct {
	handler http.Handler
}

func (logger *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	logger.handler.ServeHTTP(w, r)
	log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
}

func NewLogger(handlerToWrap http.Handler) *Logger {
	return &Logger{handlerToWrap}
}

func SetupRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/register-user", handlers.HandleUserRegistration)
	mux.HandleFunc("/api/retrieve-movies", handlers.HandleMoviesRetrievel)

	loggerAndMux := NewLogger(mux)
	return loggerAndMux
}
