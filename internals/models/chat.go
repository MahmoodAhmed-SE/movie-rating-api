package models

import "time"


type Chat struct {
	Id          int    `json:"id"`
	UserId      int    `json:"user_id"`
	MovieId     int    `json:"movie_id"`
	TextContent string `json:"text_content"`
	CreatedAt   time.Time `json:"created_at"`
}


type ChatItem struct {
	UserId      int `json:"user_id"`
	MovieId     int `json:"movie_id"`
	TextContent string `json:"text_content"`
	CreatedAt   time.Time `json:"created_at"`
}