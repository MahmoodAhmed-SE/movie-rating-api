package movie

import (
	"encoding/json"
	"log"
	"movie-rating-api-go/internals/database"
	"net/http"
)

type ChatOnMovieRequestBody struct {
	MovieId int `json:"movie_id"`
}

type ChatItem struct {
	UserId      string `json:"user_id"`
	MovieId     string `json:"movie_id"`
	TextContent string `json:"text_content"`
	CreatedAt   string `json:"created_at"`
}

// current logic: retrieve the chats of said movie_id (TODO: needs to be changed)
func ChatOnMovie(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	var reqBody ChatOnMovieRequestBody
	if err := decoder.Decode(&reqBody); err != nil {
		log.Printf("Error while decoding ChatOnMovie request body in ChatOnMovieRequestBody struct: %v", err)
		http.Error(w, "Bad Request!", http.StatusBadRequest)
		return
	}

	conn := database.GetConn()

	rows, err := conn.Query("SELECT user_id, movie_id, text_content, created_at FROM CHATS WHERE movie_id = $1;", reqBody.MovieId)

	if err != nil {
		log.Printf("Error while querying chats using select statement in ChatOnMovie: %v", err)
		http.Error(w, "Internal Server Error!", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var chats []ChatItem

	for rows.Next() {
		var chatItem ChatItem

		if err := rows.Scan(&chatItem.MovieId, &chatItem.UserId, &chatItem.TextContent, &chatItem.CreatedAt); err != nil {
			log.Printf("Error while using Scan to encode returned values of db to ChatItem sturct: %v", err)
			http.Error(w, "Internal Server Error!", http.StatusInternalServerError)
			return
		}

		chats = append(chats, chatItem)
	}

	log.Println(chats)

	encoder := json.NewEncoder(w)

	if err := encoder.Encode(chats); err != nil {
		log.Printf("Error while encoding chats to json in ChatOnMovie: %v", err)
		http.Error(w, "Internal Server Error!", http.StatusInternalServerError)
		return
	}
}
