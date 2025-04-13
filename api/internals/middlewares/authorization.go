package middlewares

import (
	"context"
	"fmt"
	"log"
	constants "movie-rating-api-go/internals"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ReqToken struct {
	Token string `json:"token"`
}

func JWTAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		headerToken := r.Header.Get("Authorization")
		if headerToken == "" {
			http.Error(w, "Unauthorized request", http.StatusUnauthorized)
			log.Printf("Error token is not present in request header.")
			r.Body.Close()
			return
		}

		_, scanErr := fmt.Sscanf(headerToken, "Bearer %s", &headerToken)
		if scanErr != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Printf("Error scanning header authorization token: %v", scanErr)
			r.Body.Close()
			return
		}

		token, err := jwt.Parse(headerToken, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(os.Getenv(constants.EnvJWTSecretKey)), nil
		})

		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			log.Printf("Error parsing token: %v", err)
			r.Body.Close()
			return
		}

		if !token.Valid {
			http.Error(w, "Bad request", http.StatusBadRequest)
			log.Println("invalid token")
			r.Body.Close()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if time.Now().Unix() > int64(claims["exp"].(float64)) {
			http.Error(w, "Bad request", http.StatusBadRequest)
			log.Printf("expired token: %v", err)
			r.Body.Close()
		}

		// Add id claim to the request context for subsequent handlers
		ctx := context.WithValue(r.Context(), constants.UserIdKey, int(claims["id"].(float64)))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
