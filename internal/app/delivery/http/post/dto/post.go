package dto

import "time"

type PostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type CreatePostResponse struct {
	Slug string `json:"post_slug"`
}
type PostResponse struct {
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Slug      string    `json:"post_slug"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
