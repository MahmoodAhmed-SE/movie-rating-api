package auth

import (
	"encoding/json"
	"log"
	"movie-rating-api-go/internals/database"
	"movie-rating-api-go/internals/psql_errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type registerationRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func UserRegistration(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var data registerationRequestBody
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Error decoding registeration request body: %v", err)
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

		http.Error(w, "Internal server error", http.StatusInternalServerError)

		log.Printf("Error inserting new user %s: %v", data.Username, err)
		return
	}

	w.WriteHeader(200)
}
