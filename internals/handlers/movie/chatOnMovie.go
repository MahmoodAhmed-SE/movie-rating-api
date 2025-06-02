package movie

import (
	"encoding/json"
	"log"
	constants "movie-rating-api-go/internals"
	"movie-rating-api-go/internals/database"
	"net/http"

	"database/sql"
)

type PostReqBodyChatOnMovie struct {
	MovieId int    `json:"movie_id"`
	Message string `json:"message"`
}

// current logic: retrieve the chats of said movie_id (TODO: needs to be changed)
func ChatOnMovie(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	reqMethod := r.Method

	if reqMethod == http.MethodPost {
		decoder := json.NewDecoder(r.Body)

		var reqBody PostReqBodyChatOnMovie

		if err := decoder.Decode(&reqBody); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			log.Printf("Error while decoding request body of PostReqBodyChatOnMovie: %v", err)
			return
		}

		conn := database.GetConn()

		// check if movie with provided movie_id exists in db
		var assumedId int
		if err := conn.QueryRow("SELECT 1 FROM MOVIES WHERE id = $1;", reqBody.MovieId).Scan(&assumedId); err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Not Found", http.StatusNotFound)
				log.Printf("Error querying movie with id %v: %v", reqBody.MovieId, err)
			} else {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Printf("Error querying movie with id %v: %v", reqBody.MovieId, err)
			}

			return
		}

		userId, ok := r.Context().Value(constants.UserIdKey).(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			log.Println("User id is not accessible in the context even though the user accessed the handler [ChatOnMovie] and reached this code. Check authorization middleware where we handle the context for subsequent handlers")
			return
		}

		if _, err := conn.Exec("INSERT INTO CHATS(id, movie_id, user_id, text_content, created_at) VALUES(DEFAULT, $1, $2, $3, DEFAULT);", reqBody.MovieId, userId, reqBody.Message); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("Error while inserting a chat into chats table: %v", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Chat is posted successfully!"))
	} else {
		log.Printf("Method Not Allowed: %s", reqMethod)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
