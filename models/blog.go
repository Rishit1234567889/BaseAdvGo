package models

type Blog struct { //2.1
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	UserID  int    `json:"user_id"`
	Created string `json:"created"`
	Updated string `json:"updated"`
}
