package helpers

import "goproject/internal/app/delivery/http/post/dto"

type ResponseWithoutDataAndError struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ResponseWithError struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   any    `json:"error"`
}

type ResponseWithData struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type MyListResponse struct {
	Slug        string `json:"list_slug"`
	Name        string `json:"name"`
	Description string `json:"description"`
	NumOfPosts  int    `json:"num_of_posts"`
}

type PostsInMyListResponse struct {
	Slug        string             `json:"list_slug"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Posts       []dto.PostResponse `json:"posts"`
}
