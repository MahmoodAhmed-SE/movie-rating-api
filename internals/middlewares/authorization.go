package middlewares

import (
	"fmt"
	"log"
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
			log.Printf("Error token is not present in request header.")
			http.Error(w, "unauthorized request", http.StatusUnauthorized)
		}

		_, scanErr := fmt.Sscanf(headerToken, "bearer %s", &headerToken)
		if scanErr != nil {
			log.Printf("Error scanning header authorization token: %v", scanErr)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		// TO-DO: make sure user exists.
		key := []byte(os.Getenv("JWT_KEY"))
		token, err := jwt.Parse(headerToken, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return key, nil
		})

		if err != nil {
			log.Printf("Error parsing token: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		if !token.Valid {
			http.Error(w, "bad request", http.StatusBadRequest)
			log.Println("invalid token")
			return
		}
		if expDate, err := token.Claims.GetExpirationTime(); err != nil || expDate.Unix() < time.Now().Unix() {
			http.Error(w, "bad request", http.StatusBadRequest)
			log.Printf("expired token: %v", err)
			return
		}

		next.ServeHTTP(w, r)
	})
}
