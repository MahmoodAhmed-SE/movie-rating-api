package auth

import (
	"encoding/json"
	"fmt"
	"log"
	constants "movie-rating-api-go/internals"
	"movie-rating-api-go/internals/services"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type loginRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	var data loginRequestBody

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		log.Printf("Error decoding body: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)

		return
	}

	user, err := services.GetUser(data.Username)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Bad request", http.StatusBadRequest)
		} else {
			log.Printf("Error while querying user: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)

		return
	}

	loadingErr := godotenv.Load()
	if loadingErr != nil {
		log.Printf("Error loading environment variables: %v", loadingErr)
		http.Error(w, "Internal server error", http.StatusInternalServerError)

		return
	}

	key := []byte(os.Getenv(constants.EnvJWTSecretKey))
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tknString, signingErr := tkn.SignedString(key)

	if signingErr != nil {
		log.Printf("Error Signing token: %v", signingErr)
		http.Error(w, "Internal server error", http.StatusInternalServerError)

		return
	}

	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tknString))
	w.WriteHeader(http.StatusOK)
}
