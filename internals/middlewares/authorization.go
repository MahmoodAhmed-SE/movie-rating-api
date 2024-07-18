package middlewares

import (
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
			return
		}

		_, scanErr := fmt.Sscanf(headerToken, "Bearer %s", &headerToken)
		if scanErr != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Printf("Error scanning header authorization token: %v", scanErr)
			return
		}

		// TO-DO: make sure user exists.
		key := []byte(os.Getenv(constants.EnvJWTSecretKey))
		token, err := jwt.Parse(headerToken, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return key, nil
		})
		log.Println(headerToken)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			log.Printf("Error parsing token: %v", err)
			return
		}

		if !token.Valid {
			http.Error(w, "Bad request", http.StatusBadRequest)
			log.Println("invalid token")
			return
		}
		if expDate, err := token.Claims.GetExpirationTime(); err != nil || expDate.Unix() < time.Now().Unix() {
			http.Error(w, "Bad request", http.StatusBadRequest)
			log.Printf("expired token: %v", err)
			return
		}

		next.ServeHTTP(w, r)
	})
}
