package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
		// Read and store the request body
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		// Restore the request body for subsequent handlers
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		decoder := json.NewDecoder(bytes.NewBuffer(bodyBytes))
		reqToken := make(map[string]interface{})

		if err := decoder.Decode(&reqToken); err != nil {
			log.Printf("Error decoding token from request body: %v", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		if token, ok := reqToken["token"].(string); !ok {
			log.Printf("Error token is not string or not present in request body.")
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		} else {
			// TO-DO: make sure user exists.
			key := []byte(os.Getenv("JWT_KEY"))
			token, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
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
		}

		next.ServeHTTP(w, r)
	})
}
