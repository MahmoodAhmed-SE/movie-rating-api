package models

import "time"

type Watchlist struct {
	Id        int       `json:"id"`
	MovieId   int       `json:"movie_id"`
	UserId    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type WatchlistResponse struct {
	MovieId   int       `json:"movie_id"`
	UserId    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
