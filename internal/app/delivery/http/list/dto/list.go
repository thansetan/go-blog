package dto

import "goproject/internal/app/delivery/http/post/dto"

type ListRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ListResponse struct {
	Slug        string              `json:"list_slug"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Posts       *[]dto.PostResponse `json:"posts,omitempty"` // use pointer so when it's empty, it'll still be rendered
	NumOfPosts  *int                `json:"num_of_posts,omitempty"`
}
