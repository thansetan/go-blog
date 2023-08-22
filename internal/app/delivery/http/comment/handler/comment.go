package commenthandler

import (
	"goproject/internal/app/delivery/http/comment/dto"
	commentusecase "goproject/internal/app/usecase/comment"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommentHandler interface {
	CreateComment(c *gin.Context)
	GetMyComments(c *gin.Context)
	GetCommentsOnAPost(c *gin.Context)
	DeleteCommentByID(c *gin.Context)
	EditMyCommentOnAPost(c *gin.Context)
}

type CommentHandlerImpl struct {
	uc commentusecase.CommentUsecase
}

func NewCommentHandler(uc commentusecase.CommentUsecase) CommentHandler {
	return &CommentHandlerImpl{
		uc: uc,
	}
}

// @CreateComment godoc
// @Summary Create a comment on a blog post
// @Description Create a comment on a blog post by providing required data
// @Tags Comment
// @Param Body body dto.CommentRequest true "the body to create a comment"
// @Param username path string true "blog owner's username"
// @Param slug path string true "post slug"
// @Param Authorization header string true "Authorization. Use 'Bearer <your-token>'"
// @Security BearerToken
// @Produce json
// @Success 201 {object} map[string]any
// @Router /blog/{username}/posts/{slug}/comments [post]
func (handler *CommentHandlerImpl) CreateComment(c *gin.Context) {
	var data dto.CommentRequest
	blogOwner := c.Param("username")
	postSlug := c.Param("slug")
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
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "username not provided",
		})
	}

	err = handler.uc.CreateComment(c, data, username, blogOwner, postSlug)
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

// @GetMyComments Godoc
// @Summary Get current user's comments
// @Description Get current user's comments on all posts
// @Tags Comment
// @Param Authorization header string true "Authorization. Use 'Bearer <your-token>'"
// @Security BearerToken
// @Produce json
// @Success 200 {objects} map[string]any
// @Router /my/comments [get]
func (handler *CommentHandlerImpl) GetMyComments(c *gin.Context) {
	username := c.GetString("username")
	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "username not provided",
		})
		return
	}

	comments, err := handler.uc.GetCommentByUsername(c, username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    comments,
	})
}

// @GetCommentOnAPost Godoc
// @Summary Get comments on a post
// @Description Get all comments on a post
// @Tags Comment
// @Param username path string true "blog owner's username"
// @Param slug path string true "post slug"
// @Produce json
// @Success 200 {objects} map[string]any
// @Router /blog/{username}/posts/{slug}/comments [get]
func (handler *CommentHandlerImpl) GetCommentsOnAPost(c *gin.Context) {
	blogOwner := c.Param("username")
	postSlug := c.Param("slug")

	comments, err := handler.uc.GetCommentByBlogOwnerAndPostSlug(c, blogOwner, postSlug)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    comments,
	})
}

// @DeleteCommentByID Godoc
// @Summary Delete comment by ID
// @Description Delete comment by ID
// @Tags Comment
// @Param username path string true "blog owner's username"
// @Param slug path string true "post slug"
// @Param comment_id path int true "comment ID"
// @Param Authorization header string true "Authorization. Use 'Bearer <your-token>'"
// @Security BearerToken
// @Produce json
// @Success 200 {objects} map[string]any
// @Router /blog/{username}/posts/{slug}/comments/{comment_id} [delete]
func (handler *CommentHandlerImpl) DeleteCommentByID(c *gin.Context) {
	blogOwner := c.Param("username")
	postSlug := c.Param("slug")
	commentID := c.Param("comment_id")
	username := c.GetString("username")

	err := handler.uc.DeleteCommentOnAPosst(c, username, blogOwner, postSlug, commentID)
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

// @EditMyCommentOnAPost Godoc
// @Summary Edit user's comment on a post
// @Description Edit user's comment on a post by providing comment ID
// @Tags Comment
// @Param username path string true "blog owner's username"
// @Param slug path string true "post slug"
// @Param Body body dto.CommentRequest true "body required to modify comment"
// @Param comment_id path int true "comment ID"
// @Param Authorization header string true "Authorization. Use 'Bearer <your-token>'"
// @Security BearerToken
// @Produce json
// @Success 200 {objects} map[string]any
// @Router /blog/{username}/posts/{slug}/comments/{comment_id} [put]
func (handler *CommentHandlerImpl) EditMyCommentOnAPost(c *gin.Context) {
	var data dto.CommentRequest
	blogOwner := c.Param("username")
	postSlug := c.Param("slug")
	commentID := c.Param("comment_id")
	username := c.GetString("username")

	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "username not provided",
		})
	}

	err = handler.uc.UpdateCommentOnAPost(c, username, blogOwner, postSlug, commentID, data)
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