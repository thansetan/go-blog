package dto

import "time"

type CommentRequest struct {
	Comment string `json:"comment"`
}

type CommentResponse struct {
	ID        uint      `json:"comment_id"`
	PostURL   string    `json:"post_url,omitempty"`
	Commenter string    `json:"commenter,omitempty"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
