package dto

import "time"

type PostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type PostResponse struct {
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Slug      string    `json:"title-slug"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
