package bloghandler

import (
	"fmt"
	"goproject/internal/app/delivery/http/blog/dto"
	blogusecase "goproject/internal/app/usecase/blog"
	"goproject/internal/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BlogHandler interface {
	UpdateBlogData(c *gin.Context)
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
// @Security BearerToken
// @Produce json
// @Success 200 {object} map[string]any
// @Router /blog/my [put]
func (handler *BlogHandlerImpl) UpdateBlogData(c *gin.Context) {
	var blog dto.UpdateBlogRequest
	username := c.GetString("username")

	if username == "" {
		helpers.ResponseBuilder(c, http.StatusUnauthorized, "update blog data", "you're not allowed to access this path", nil)
		return
	}

	err := c.ShouldBindJSON(&blog)
	if err != nil {
		helpers.ResponseBuilder(c, http.StatusBadRequest, "update blog data", helpers.ValidationError(err), nil)
		return
	}

	err = handler.uc.UpdateBlogData(c, username, blog)
	if err != nil {
		helpers.ResponseBuilder(c, http.StatusInternalServerError, "update bog data", err.Error(), nil)
		return
	}

	helpers.ResponseBuilder(c, http.StatusOK, "update blog data", nil, nil)
}

// GetMyBlog godoc
// @Summary Get my blog information
// @Description Get information about my blog
// @Tags Blog
// @Security BearerToken
// @Produce json
// @Success 200 {object} map[string]any
// @Router /blog/my [get]
func (handler *BlogHandlerImpl) GetMyBlog(c *gin.Context) {
	username := c.GetString("username")
	if username == "" {
		helpers.ResponseBuilder(c, http.StatusUnauthorized, "get my blog", "you're not allowed to access this path", nil)
		return
	}

	blog, err := handler.uc.GetBlogByOwner(c, username)
	if err != nil {
		helpers.ResponseBuilder(c, http.StatusInternalServerError, "get my blog", err.Error(), nil)
		return
	}

	helpers.ResponseBuilder(c, http.StatusOK, "get my blog", nil, blog)
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
		helpers.ResponseBuilder(c, http.StatusInternalServerError, fmt.Sprintf("get %s's blog", owner), err.Error(), nil)
		return
	}

	helpers.ResponseBuilder(c, http.StatusOK, fmt.Sprintf("get %s's blog", owner), nil, blog)
}
