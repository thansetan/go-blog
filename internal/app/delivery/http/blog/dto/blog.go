package dto

type UpdateBlogRequest struct {
	NewName string `json:"new_name"`
}

type BlogResponse struct {
	Name  string `json:"name"`
	Owner string `json:"blog_owner"`
}
