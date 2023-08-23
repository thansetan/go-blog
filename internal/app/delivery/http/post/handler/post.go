package posthandler

import (
	"goproject/internal/app/delivery/http/post/dto"
	postusecase "goproject/internal/app/usecase/post"
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

type PostHandlerImpl struct {
	uc postusecase.PostUsecase
}

func NewPostHandler(uc postusecase.PostUsecase) PostHandler {
	return &PostHandlerImpl{
		uc: uc,
	}
}

// @CreateNewPost godoc
// @Summary Create a new blog post
// @Description Create a new blog post by providing required data
// @Tags Post
// @Param Body body dto.PostRequest true "the body to create a new post"
// @Security BearerToken
// @Produce json
// @Success 201 {object} map[string]any
// @Router /blog/my/posts [post]
func (handler *PostHandlerImpl) CreateNewPost(c *gin.Context) {
	var data dto.PostRequest
	username := c.GetString("username")

	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "username not provided",
		})
		return
	}

	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
		})
		return
	}

	err = handler.uc.CreateNewPost(c, username, data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
	})
}

// @GetMyPosts godoc
// @Summary Get all user's blog posts
// @Description Get all user's blog posts by providing JWT auth
// @Tags Post
// @Security BearerToken
// @Produce json
// @Success 200 {object} map[string]any
// @Router /blog/my/posts [GET]
func (handler *PostHandlerImpl) GetAllMyBlogPosts(c *gin.Context) {
	username := c.GetString("username")

	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "username not provided",
		})
		return
	}

	posts, err := handler.uc.GetPostsByBlogOwner(c, username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    posts,
	})
}

// @GetUsersBlogPosts godoc
// @Summary Get all blog posts of a user
// @Description Get all user's blog posts by providing username
// @Tags Post
// @Param username path string true "Username of the user"
// @Produce json
// @Success 200 {object} map[string]any
// @Router /blog/{username}/posts [GET]
func (handler *PostHandlerImpl) GetPostsByBlogOwner(c *gin.Context) {
	username := c.Param("username")

	posts, err := handler.uc.GetPostsByBlogOwner(c, username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    posts,
	})
}

// @GetUserPostBySlug godoc
// @Summary Get a user post by slug
// @Description Get a user post by providing their username and the post slug
// @Tags Post
// @Param username path string true "Username of the user"
// @Param post_slug path string true "Slug of the post"
// @Produce json
// @Success 200 {object} map[string]any
// @Router /blog/{username}/posts/{post_slug} [GET]
func (handler *PostHandlerImpl) GetPostBySlug(c *gin.Context) {
	slug := c.Param("post_slug")
	username := c.Param("username")
	post, err := handler.uc.GetPostBySlug(c, username, slug)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    post,
	})
}

// @UpdateMyPostBySlug godoc
// @Summary Update user post by slug
// @Description Update user blog post by providing the post slug
// @Tags Post
// @Security BearerToken
// @Param Body body dto.PostRequest true "the body to create a new post"
// @Param post_slug path string true "Slug of the post"
// @Produce json
// @Success 200 {object} map[string]any
// @Router /blog/my/posts/{post_slug} [PUT]
func (handler *PostHandlerImpl) UpdateMyPostBySlug(c *gin.Context) {
	var data dto.PostRequest
	slug := c.Param("post_slug")
	username := c.GetString("username")

	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "username not provided",
		})
		return
	}

	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	err = handler.uc.UpdatePostBySlug(c, data, username, slug)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// @DeleteMyPostBySlug godoc
// @Summary Delete user post by slug
// @Description Delete user blog post by providing the post slug
// @Tags Post
// @Security BearerToken
// @Param post_slug path string true "Slug of the post"
// @Produce json
// @Success 200 {object} map[string]any
// @Router /blog/my/posts/{post_slug} [DELETE]
func (handler *PostHandlerImpl) DeleteMyPostBySlug(c *gin.Context) {
	slug := c.Param("post_slug")
	username := c.GetString("username")

	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "username not provided",
		})
		return
	}
	err := handler.uc.DeletePostBySlug(c, username, slug)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
