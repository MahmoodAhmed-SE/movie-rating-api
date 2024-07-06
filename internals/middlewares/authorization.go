package middlewares

import (
	"encoding/json"
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
		decoder := json.NewDecoder(r.Body)
		var reqToken ReqToken
		if err := decoder.Decode(&reqToken); err != nil {
			log.Printf("Error decoding request body: %v", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		log.Println(reqToken.Token)
		// TO-DO: make sure user exists.
		key := []byte(os.Getenv("JWT_KEY"))
		token, err := jwt.Parse(reqToken.Token, func(t *jwt.Token) (interface{}, error) {
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
			http.Error(w, "invalid jwt", http.StatusBadRequest)
			return
		}
		if expDate, err := token.Claims.GetExpirationTime(); err != nil || expDate.Unix() < time.Now().Unix() {
			http.Error(w, "invalid jwt", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
