package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"movie-rating-api-go/internals/database"
	"movie-rating-api-go/internals/psql_errors"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type registerationRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func HandleUserRegistration(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var data registerationRequestBody
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		log.Fatal(err)
		return
	}

	// saving user data into postgresdb
	hashedPasswordBytes, hashErr := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)

	if hashErr != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Error generating hash: %v", hashErr)
		return
	}

	conn := database.GetConn()
	insertionRow := conn.QueryRow("INSERT INTO users(username, password) VALUES($1, $2) RETURNING id;", data.Username, string(hashedPasswordBytes))

	// var userId int
	var userId int
	if err := insertionRow.Scan(&userId); err != nil {
		if err.Error() == psql_errors.UniqueConstraintViolation {
			http.Error(w, "Username already exists!", http.StatusConflict)
			return
		}

		log.Printf("Error inserting new user "+data.Username+": %v", err)

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var (
		key []byte
		t   *jwt.Token
	)

	key = []byte(os.Getenv("JWT_KEY"))
	t = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  fmt.Sprint(userId),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	s, err := t.SignedString(key)
	if err != nil {
		log.Printf("Error returning signed jwt token: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s))
	w.WriteHeader(200)
	w.Write([]byte(s))
	log.Printf("Insertion Id:%d\n%s", userId, s)
}
