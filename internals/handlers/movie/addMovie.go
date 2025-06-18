package movie

import (
	"encoding/json"
	"log"
	"movie-rating-api-go/internals/database"
	"movie-rating-api-go/internals/services"
	"net/http"
)

type httpRequst struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func AddMovie(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	var reqBody httpRequst
	if err := decoder.Decode(&reqBody); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println(err)
		return
	}

	// TO-DO: add all information of a movie to the description vec (name, directors, etc..)
	// TO-DO: request body validation

	description_vec, err := services.GetGrpcEmbeddingResp(reqBody.Description)
	if err != nil {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		log.Println(err)
		return
	}

	dbConn := database.GetConn()

	_, err = dbConn.Exec("INSERT INTO movies(id, name, description, description_vec) VALUES(DEFAULT, $1, $2, $3);", reqBody.Name, reqBody.Description, description_vec)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	log.Printf("Added movie: %s", reqBody.Name)
	type httpResp struct {
		Message string `json:"message"`
	}

	userResp, err := json.Marshal(&httpResp{
		Message: "movie successfully added!",
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Write(userResp)
}
