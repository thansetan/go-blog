package dto

type UpdateBlogRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type BlogResponse struct {
	Name        string `json:"blog_name"`
	Description string `json:"blog_description"`
	Owner       string `json:"blog_owner"`
}
