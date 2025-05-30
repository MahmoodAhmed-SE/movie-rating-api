package models

type Chat struct {
	Id          int    `json:"id"`
	UserId      int    `json:"user_id"`
	MovieId     int    `json:"movie_id"`
	TextContent string `json:"text_content"`
	CreatedAt   string `json:"created_at"`
}


type ChatItem struct {
	UserId      string `json:"user_id"`
	MovieId     string `json:"movie_id"`
	TextContent string `json:"text_content"`
	CreatedAt   string `json:"created_at"`
}