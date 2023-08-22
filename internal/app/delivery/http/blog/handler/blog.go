package bloghandler

import (
	"goproject/internal/app/delivery/http/blog/dto"
	blogusecase "goproject/internal/app/usecase/blog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BlogHandler interface {
	UpdateMyBlogName(c *gin.Context)
	GetMyBlog(c *gin.Context)
	GetBlogByOwner(c *gin.Context)
}

type BlogHandlerImpl struct {
	uc blogusecase.BlogUsecase
}

func NewBlogHandler(uc blogusecase.BlogUsecase) BlogHandler {
	return &BlogHandlerImpl{
		uc: uc,
	}
}

// UpdateMyBlog godoc
// @Summary Change user's blog name
// @Description Change user's blog name by providing required data
// @Tags Blog
// @Param Body body dto.UpdateBlogRequest true "the body to change user's blog name"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Success 200 {object} map[string]any
// @Router /blog/my [put]
func (handler *BlogHandlerImpl) UpdateMyBlogName(c *gin.Context) {
	var blog dto.UpdateBlogRequest
	username := c.GetString("username")

	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "username not provided",
		})
		return
	}

	err := c.ShouldBindJSON(&blog)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
		})
		return
	}

	err = handler.uc.UpdateBlogName(c, username, blog)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// GetMyBlog godoc
// @Summary Get my blog information
// @Description Get information about my blog
// @Tags Blog
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Security BearerToken
// @Produce json
// @Success 200 {object} map[string]any
// @Router /blog/my [get]
func (handler *BlogHandlerImpl) GetMyBlog(c *gin.Context) {
	username := c.GetString("username")
	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "username not provided",
		})
		return
	}

	blog, err := handler.uc.GetBlogByOwner(c, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    blog,
	})
}

// GetUserBlog godoc
// @Summary Get user blog information
// @Description Get information about user blog
// @Tags Blog
// @Param username path string true "Username of the user"
// @Produce json
// @Success 200 {object} map[string]any
// @Router /blog/{username} [get]
func (handler *BlogHandlerImpl) GetBlogByOwner(c *gin.Context) {
	owner := c.Param("username")

	blog, err := handler.uc.GetBlogByOwner(c, owner)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    blog,
	})
}
