package models

type Chat struct {
	Id          int    `json:"id"`
	UserId      int    `json:"user_id"`
	MovieId     int    `json:"movie_id"`
	TextContent string `json:"text_content"`
	CreatedAt   string `json:"created_at"`
}
