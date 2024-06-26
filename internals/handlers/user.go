package handlers

import (
	"encoding/json"
	"log"
	"movie-rating-api-go/internals/database"
	"movie-rating-api-go/internals/psql_errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type RequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func HandleUserRegistration(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var data RequestBody
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		log.Fatal(err)
		return
	}

	// saving user data into postgresdb
	hashedPasswordBytes, hashErr := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)

	if hashErr != nil {
		switch hashErr.(type) {
		case bcrypt.InvalidCostError:
			http.Error(w, "Invalid cost parameter", http.StatusBadRequest)
		case bcrypt.InvalidHashPrefixError:
			http.Error(w, "Invalid hash prefix", http.StatusBadRequest)
		case bcrypt.HashVersionTooNewError:
			http.Error(w, "Hash version too new", http.StatusBadRequest)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
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

	log.Printf("Insertion: \nId:%dn", userId)
}
