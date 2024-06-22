package main

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type RequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func handleUserRegistration(w http.ResponseWriter, r *http.Request) {
	var data RequestBody
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&data); err != nil {
		log.Fatal(err)
	}

	defer r.Body.Close()

	// saving user data into postgresdb
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)

	if err != nil {
		switch err.(type) {
		case bcrypt.InvalidCostError:
			http.Error(w, "Invalid cost parameter", http.StatusBadRequest)
		case bcrypt.InvalidHashPrefixError:
			http.Error(w, "Invalid hash prefix", http.StatusBadRequest)
		case bcrypt.HashVersionTooNewError:
			http.Error(w, "Hash version too new", http.StatusBadRequest)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		log.Printf("Error generating hash: %v", err)
		return
	}

}

func main() {

	http.HandleFunc("/api/", handleUserRegistration)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}

}
