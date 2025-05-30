package movie

import (
	"encoding/json"
	"log"
	"movie-rating-api-go/internals/database"
	"net/http"
	"github.com/gorilla/mux"
)


type ChatItem struct {
	UserId      string `json:"user_id"`
	MovieId     string `json:"movie_id"`
	TextContent string `json:"text_content"`
	CreatedAt   string `json:"created_at"`
}

// current logic: retrieve the chats of said movie_id (TODO: needs to be changed)
func ChatOnMovie(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	reqMethod := r.Method

	if reqMethod == http.MethodGet {
		vars := mux.Vars(r)
		movieId := vars["movieId"]


		conn := database.GetConn()

		rows, err := conn.Query("SELECT user_id, movie_id, text_content, created_at FROM CHATS WHERE movie_id = $1;", movieId)

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
	} else {
		log.Printf("Method Not Allowed: %s", reqMethod)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
