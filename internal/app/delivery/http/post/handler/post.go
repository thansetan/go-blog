package posthandler

import (
	"fmt"
	"goproject/internal/app/delivery/http/post/dto"
	postusecase "goproject/internal/app/usecase/post"
	"goproject/internal/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PostHandler interface {
	CreateNewPost(c *gin.Context)
	GetPostsByBlogOwner(c *gin.Context)
	GetAllMyBlogPosts(c *gin.Context)
	GetPostBySlug(c *gin.Context)
	UpdateMyPostBySlug(c *gin.Context)
	DeleteMyPostBySlug(c *gin.Context)
}

type postHandlerImpl struct {
	uc postusecase.PostUsecase
}

func NewPostHandler(uc postusecase.PostUsecase) PostHandler {
	return &postHandlerImpl{
		uc: uc,
	}
}

//	@CreateNewPost	godoc
//	@Summary		Create a new blog post
//	@Description	Create a new post on current user's blog.
//	@Description	Upon successful creation, it will return the newly created post's slug
//	@Tags			Post
//	@Param			Body	body	dto.PostRequest	true	"data required to create a new post"
//	@Security		BearerToken
//	@Produce		json
//	@Success		201	{object}	helpers.ResponseWithData{data=dto.CreatePostResponse}
//	@Router			/blog/my/posts [post]
func (handler *postHandlerImpl) CreateNewPost(c *gin.Context) {
	var data dto.PostRequest
	username := c.GetString("username")

	if username == "" {
		helpers.ResponseBuilder(c, http.StatusUnauthorized, "create post", "you're not allowed to access this path", nil)
		return
	}

	err := c.ShouldBindJSON(&data)
	if err != nil {
		helpers.ResponseBuilder(c, http.StatusBadRequest, "create post", helpers.ValidationError(err), nil)
		return
	}

	resp, ucErr := handler.uc.CreateNewPost(c, username, data)
	if ucErr != nil {
		helpers.ResponseBuilder(c, ucErr.Code, "create post", ucErr.String(), nil)
		return
	}

	helpers.ResponseBuilder(c, http.StatusCreated, "create post", nil, resp)
}

//	@GetMyPosts		godoc
//	@Summary		Get all current user's blog posts
//	@Description	Get all current user's blog posts. When there are no posts, it will return an empty array ([]).
//	@Tags			Post
//	@Security		BearerToken
//	@Produce		json
//	@Success		200	{object}	helpers.ResponseWithData{data=[]dto.PostResponse}
//	@Router			/blog/my/posts [GET]
func (handler *postHandlerImpl) GetAllMyBlogPosts(c *gin.Context) {
	username := c.GetString("username")

	if username == "" {
		helpers.ResponseBuilder(c, http.StatusUnauthorized, "get my posts", "you're not allowed to access this path", nil)
		return
	}

	posts, err := handler.uc.GetPostsByBlogOwner(c, username)
	if err != nil {
		helpers.ResponseBuilder(c, err.Code, "get my posts", err.String(), nil)
		return
	}

	helpers.ResponseBuilder(c, http.StatusOK, "get my posts", nil, posts)
}

//	@GetUsersBlogPosts	godoc
//	@Summary			Get a user's blog posts
//	@Description		Get user's blog posts by providing their username.
//	@Tags				Post
//	@Param				username	path	string	true	"user's username"
//	@Produce			json
//	@Success			200	{object}	helpers.ResponseWithData{data=[]dto.PostResponse}
//	@Failure			404	{object}	helpers.ResponseWithError
//	@Router				/blog/{username}/posts [GET]
func (handler *postHandlerImpl) GetPostsByBlogOwner(c *gin.Context) {
	username := c.Param("username")

	posts, err := handler.uc.GetPostsByBlogOwner(c, username)
	if err != nil {
		helpers.ResponseBuilder(c, err.Code, fmt.Sprintf("get %s's posts", username), err.String(), nil)
		return
	}

	helpers.ResponseBuilder(c, http.StatusOK, fmt.Sprintf("get %s's posts", username), nil, posts)
}

//	@GetUserPostBySlug	godoc
//	@Summary			Get a specific post
//	@Description		Get a specific post by providing their username and the post's slug.
//	@Tags				Post
//	@Param				username	path	string	true	"user's username"
//	@Param				post_slug	path	string	true	"post's slug"
//	@Produce			json
//	@Success			200	{object}	helpers.ResponseWithData{data=dto.PostResponse}
//	@Failure			404	{object}	helpers.ResponseWithError
//	@Router				/blog/{username}/posts/{post_slug} [GET]
func (handler *postHandlerImpl) GetPostBySlug(c *gin.Context) {
	slug := c.Param("post_slug")
	username := c.Param("username")
	post, err := handler.uc.GetPostBySlug(c, username, slug)
	if err != nil {
		helpers.ResponseBuilder(c, err.Code, "get post", err.String(), nil)
		return
	}

	helpers.ResponseBuilder(c, http.StatusOK, "get post", nil, post)
}

//	@UpdateMyPostBySlug	godoc
//	@Summary			Update/modify current user's post
//	@Description		Update/modify current user's blog post by providing the post's slug.
//	@Tags				Post
//	@Security			BearerToken
//	@Param				Body		body	dto.PostRequest	true	"data required to update/modify a post"
//	@Param				post_slug	path	string			true	"post's slug"
//	@Produce			json
//	@Success			200	{object}	helpers.ResponseWithoutDataAndError
//	@Failure			404	{object}	helpers.ResponseWithError
//	@Router				/blog/my/posts/{post_slug} [PUT]
func (handler *postHandlerImpl) UpdateMyPostBySlug(c *gin.Context) {
	var data dto.PostRequest
	slug := c.Param("post_slug")
	username := c.GetString("username")

	if username == "" {
		helpers.ResponseBuilder(c, http.StatusUnauthorized, "update post", "you're not allowed to access this path", nil)
		return
	}

	err := c.ShouldBindJSON(&data)
	if err != nil {
		helpers.ResponseBuilder(c, http.StatusBadRequest, "update post", helpers.ValidationError(err), nil)
		return
	}

	ucErr := handler.uc.UpdatePostBySlug(c, data, username, slug)
	if ucErr != nil {
		helpers.ResponseBuilder(c, ucErr.Code, "update post", ucErr.String(), nil)
		return
	}

	helpers.ResponseBuilder(c, http.StatusOK, "update post", nil, nil)
}

//	@DeleteMyPostBySlug	godoc
//	@Summary			Delete current user's post
//	@Description		Delete current user's blog post by providing the post slug.
//	@Description		When a post is deleted, all comments on the post will also be deleted.
//	@Description		Also, if the post is part of someone's lists, the post in that list will also be deleted.
//	@Tags				Post
//	@Security			BearerToken
//	@Param				post_slug	path	string	true	"post's slug"
//	@Produce			json
//	@Success			200	{object}	helpers.ResponseWithoutDataAndError
//	@Failure			404	{object}	helpers.ResponseWithError
//	@Router				/blog/my/posts/{post_slug} [DELETE]
func (handler *postHandlerImpl) DeleteMyPostBySlug(c *gin.Context) {
	slug := c.Param("post_slug")
	username := c.GetString("username")

	if username == "" {
		helpers.ResponseBuilder(c, http.StatusUnauthorized, "delete post", "you're not allowed to access this path", nil)
		return
	}
	err := handler.uc.DeletePostBySlug(c, username, slug)
	if err != nil {
		helpers.ResponseBuilder(c, err.Code, "delete post", err.String(), nil)
		return
	}

	helpers.ResponseBuilder(c, http.StatusOK, "delete post", nil, nil)
}
