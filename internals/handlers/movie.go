package handlers

import (
	"encoding/json"
	"log"
	"net/http"
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
		log.Printf("Error decoding movies retrievel api request body: %v", err)
		http.Error(w, "bad request", http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}
