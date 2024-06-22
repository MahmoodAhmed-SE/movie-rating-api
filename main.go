package main

import (
	"encoding/json"
	"log"
	"net/http"
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

}

func main() {

	http.HandleFunc("/api/", handleUserRegistration)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}

}
