package movie

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"movie-rating-api-go/internals/database"
	"movie-rating-api-go/internals/go_protobuffs"
	"net/http"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type httpRequst struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func floatSliceToVectorLiteral(vec []float32) string {
	strVals := make([]string, len(vec))
	for i, v := range vec {
		strVals[i] = fmt.Sprintf("%f", v)
	}
	return fmt.Sprintf("[%s]", strings.Join(strVals, ","))
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

	// TO-DO: request body validation

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		log.Println(err)
		return
	}
	defer conn.Close()

	client := go_protobuffs.NewEmbedderClient(conn)

	started := time.Now().UnixMilli()
	resp, err := client.Embed(context.Background(), &go_protobuffs.EmbeddingRequest{
		Text: reqBody.Description,
	})
	log.Println("Time taken:", time.Now().UnixMilli()-started)

	if err != nil {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		log.Println(err)
		return
	}

	dbConn := database.GetConn()

	description_vec := floatSliceToVectorLiteral(resp.GetEmbedding())

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
