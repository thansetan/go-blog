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
//	@Summary		Change current user's blog information
//	@Description	Change current user's blog name and description.
//	@Tags			Blog
//	@Param			Body	body	dto.UpdateBlogRequest	true	"data required to change user's blog information"
//	@Security		BearerToken
//	@Produce		json
//	@Success		200	{object}	helpers.ResponseWithoutDataAndError
//	@Router			/blog/my [put]
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

	ucErr := handler.uc.UpdateBlogData(c, username, blog)
	if ucErr != nil {
		helpers.ResponseBuilder(c, ucErr.Code, "update bog data", ucErr.String(), nil)
		return
	}

	helpers.ResponseBuilder(c, http.StatusOK, "update blog data", nil, nil)
}

// GetMyBlog godoc
//	@Summary		Get current user's blog information
//	@Description	Get current user's blog information (name, description, number of posts).
//	@Tags			Blog
//	@Security		BearerToken
//	@Produce		json
//	@Success		200	{object}	helpers.ResponseWithData{data=dto.BlogResponse}
//	@Router			/blog/my [get]
func (handler *BlogHandlerImpl) GetMyBlog(c *gin.Context) {
	username := c.GetString("username")
	if username == "" {
		helpers.ResponseBuilder(c, http.StatusUnauthorized, "get my blog", "you're not allowed to access this path", nil)
		return
	}

	blog, err := handler.uc.GetBlogByOwner(c, username)
	if err != nil {
		helpers.ResponseBuilder(c, err.Code, "get my blog", err.String(), nil)
		return
	}

	helpers.ResponseBuilder(c, http.StatusOK, "get my blog", nil, blog)
}

// GetUserBlog godoc
//	@Summary		Get user's blog information
//	@Description	Get user's blog information (name, description, number of posts) by providing their username.
//	@Tags			Blog
//	@Param			username	path	string	true	"Username of the user"
//	@Produce		json
//	@Success		200	{object}	helpers.ResponseWithData{data=dto.BlogResponse}
//	@Failure		404	{object}	helpers.ResponseWithError{error=string}
//	@Router			/blog/{username} [get]
func (handler *BlogHandlerImpl) GetBlogByOwner(c *gin.Context) {
	owner := c.Param("username")

	blog, err := handler.uc.GetBlogByOwner(c, owner)
	if err != nil {
		helpers.ResponseBuilder(c, err.Code, fmt.Sprintf("get %s's blog", owner), err.String(), nil)
		return
	}

	helpers.ResponseBuilder(c, http.StatusOK, fmt.Sprintf("get %s's blog", owner), nil, blog)
}
