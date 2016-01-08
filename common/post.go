package common

import "time"

type NewPost struct {
	Message string `json:"message"`
	Title   string `json:"title"`
}

type Post struct {
	Message   string    `json:"message"`
	Title     string    `json:"title"`
	Username  string    `json:"username"`
	Timestamp time.Time `json:"timestamp"`
}

type Posts []Post
