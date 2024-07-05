package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type movieRetrievelRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

func HandleMoviesRetrievel(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var data movieRetrievelRequestBody

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "bad request", http.StatusBadRequest)
	}
	log.Println(data.Token)
	// TO-DO: make sure user exists.
	key := []byte(os.Getenv("JWT_KEY"))
	token, err := jwt.Parse(string(data.Token), jwt.Keyfunc(func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return key, nil
	}))

	if err != nil {
		log.Printf("Error parsing token: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if !token.Valid {
		http.Error(w, "invalid jwt", http.StatusBadRequest)
		return
	}
	if expDate, err := token.Claims.GetExpirationTime(); err != nil {
		if expDate.Unix() < time.Now().Unix() {
			http.Error(w, "invalid jwt", http.StatusBadRequest)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
