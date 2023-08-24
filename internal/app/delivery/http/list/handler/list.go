package listhandler

import (
	"goproject/internal/app/delivery/http/list/dto"
	listusecase "goproject/internal/app/usecase/list"
	"goproject/internal/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ListHandler interface {
	CreateNewList(c *gin.Context)
	AddPostToMyList(c *gin.Context)
	GetPostsInMyListBySlug(c *gin.Context)
	GetMyLists(c *gin.Context)
	UpdateMyListInformationBySlug(c *gin.Context)
	RemovePostFromMyList(c *gin.Context)
	DeleteMyListBySlug(c *gin.Context)
}

type ListHandlerImpl struct {
	uc listusecase.ListUsecase
}

func NewListHandler(uc listusecase.ListUsecase) ListHandler {
	return &ListHandlerImpl{
		uc: uc,
	}
}

// @CreateNewList godoc
// @Summary Create a new list for user
// @Description Create a new list for user by providing required data
// @Tags List
// @Param Body body dto.ListRequest true "body required to create a new list"
// @Security BearerToken
// @Produce json
// @Success 201 {object} map[string]any
// @Router /lists/my [post]
func (handler *ListHandlerImpl) CreateNewList(c *gin.Context) {
	var data dto.ListRequest
	username := c.GetString("username")

	if username == "" {
		helpers.ResponseBuilder(c, http.StatusUnauthorized, "create list", "you're not allowed to access this path", nil)
		return
	}

	err := c.ShouldBindJSON(&data)
	if err != nil {
		helpers.ResponseBuilder(c, http.StatusBadRequest, "create list", helpers.ValidationError(err), nil)
		return
	}

	id, err := handler.uc.CreateNewList(c, data, username)
	if err != nil {
		helpers.ResponseBuilder(c, http.StatusInternalServerError, "create list", err.Error(), nil)
		return
	}

	helpers.ResponseBuilder(c, http.StatusCreated, "create list", nil, id)
}

// @AddPostToMyList godoc
// @Summary Add post to my list
// @Description Add post to user's list by providing required data
// @Tags List
// @Param username path string true "blog owner's username"
// @Param post_slug path string true "post slug"
// @Param list_slug path string true "list slug you want to add this post to"
// @Security BearerToken
// @Produce json
// @Success 200 {object} map[string]any
// @Router /blog/{username}/posts/{post_slug}/save/{list_slug} [post]
func (handler *ListHandlerImpl) AddPostToMyList(c *gin.Context) {
	username := c.GetString("username")
	postSlug := c.Param("post_slug")
	blogOwner := c.Param("username")
	listSlug := c.Param("list_slug")

	if username == "" {
		helpers.ResponseBuilder(c, http.StatusUnauthorized, "add post to list", "you're not allowed to access this path", nil)
		return
	}

	err := handler.uc.AddPostToMyList(c, listSlug, username, blogOwner, postSlug)
	if err != nil {
		helpers.ResponseBuilder(c, http.StatusInternalServerError, "add post to list", err.Error(), nil)
		return
	}

	helpers.ResponseBuilder(c, http.StatusOK, "add post to list", nil, nil)
}

// @GetPostsInMyList godoc
// @Summary Get posts in my list
// @Description Get posts in my list by providing required data
// @Tags List
// @Param list_slug path string true "list slug you want to get"
// @Security BearerToken
// @Produce json
// @Success 200 {object} map[string]any
// @Router /lists/my/{list_slug} [get]
func (handler *ListHandlerImpl) GetPostsInMyListBySlug(c *gin.Context) {
	username := c.GetString("username")
	listSlug := c.Param("list_slug")

	if username == "" {
		helpers.ResponseBuilder(c, http.StatusUnauthorized, "get posts in list", "you're not allowed to access this path", nil)
		return
	}

	list, err := handler.uc.GetPostsInAListBySlug(c, listSlug, username)
	if err != nil {
		helpers.ResponseBuilder(c, http.StatusInternalServerError, "get posts in list", err.Error(), nil)
		return
	}

	helpers.ResponseBuilder(c, http.StatusOK, "get posts in list", nil, list)
}

// @GetMyLists godoc
// @Summary Get current user's lists
// @Description Get current user's lists by providing required data
// @Tags List
// @Security BearerToken
// @Produce json
// @Success 200 {object} map[string]any
// @Router /lists/my [get]
func (handler *ListHandlerImpl) GetMyLists(c *gin.Context) {
	username := c.GetString("username")

	if username == "" {
		helpers.ResponseBuilder(c, http.StatusUnauthorized, "get lists", "you're not allowed to access this path", nil)
		return
	}

	lists, err := handler.uc.GetMyLists(c, username)
	if err != nil {
		helpers.ResponseBuilder(c, http.StatusInternalServerError, "get lists", err.Error(), nil)
		return
	}

	helpers.ResponseBuilder(c, http.StatusOK, "get lists", nil, lists)
}

// @UpdateMyListInformation godoc
// @Summary Update current user's list by ID
// @Description Update current user's list information (name and description) by providing the list ID
// @Tags List
// @Param list_slug path string true "list slug you want to edit"
// @Param body body dto.ListRequest strue "body to update"
// @Security BearerToken
// @Produce json
// @Success 200 {object} map[string]any
// @Router /lists/my/{list_slug} [PUT]
func (handler *ListHandlerImpl) UpdateMyListInformationBySlug(c *gin.Context) {
	var data dto.ListRequest
	username := c.GetString("username")
	listSlug := c.Param("list_slug")

	if username == "" {
		helpers.ResponseBuilder(c, http.StatusUnauthorized, "update list", "you're not allowed to access this path", nil)
		return
	}

	err := c.ShouldBindJSON(&data)
	if err != nil {
		helpers.ResponseBuilder(c, http.StatusBadRequest, "update list", helpers.ValidationError(err), nil)
		return
	}

	err = handler.uc.UpdateListInformation(c, data, username, listSlug)
	if err != nil {
		helpers.ResponseBuilder(c, http.StatusInternalServerError, "update list", err.Error(), nil)
		return
	}

	helpers.ResponseBuilder(c, http.StatusOK, "update list", nil, nil)
}

// @RemovePostFromUserList godoc
// @Summary Remove post from current user's list
// @Description Remove post from current user's list by providing the list ID and post slug
// @Tags List
// @Param list_slug path string true "list slug you want to remove post from"
// @Param post_slug path string true "post slug you want to delete"
// @Security BearerToken
// @Produce json
// @Success 200 {object} map[string]any
// @Router /lists/my/{list_slug}/{post_slug} [DELETE]
func (handler *ListHandlerImpl) RemovePostFromMyList(c *gin.Context) {
	username := c.GetString("username")
	listSlug := c.Param("list_slug")
	postSlug := c.Param("post_slug")

	if username == "" {
		helpers.ResponseBuilder(c, http.StatusUnauthorized, "remove post from list", "you're not allowed to access this path", nil)
		return
	}

	err := handler.uc.RemovePostFromList(c, username, postSlug, listSlug)
	if err != nil {
		helpers.ResponseBuilder(c, http.StatusInternalServerError, "remove post from list", err.Error(), nil)
		return
	}

	helpers.ResponseBuilder(c, http.StatusOK, "remove post from list", nil, nil)
}

// @DeleteMyList godoc
// @Summary Delete current user's list by ID
// @Description Delete current user's list by providing the list ID
// @Tags List
// @Param list_slug path string true "list slug you want to remove"
// @Security BearerToken
// @Produce json
// @Success 200 {object} map[string]any
// @Router /lists/my/{list_slug} [DELETE]
func (handler *ListHandlerImpl) DeleteMyListBySlug(c *gin.Context) {
	username := c.GetString("username")
	listSlug := c.Param("list_slug")

	if username == "" {
		helpers.ResponseBuilder(c, http.StatusUnauthorized, "delete list", "you're not allowed to access this path", nil)
		return
	}

	err := handler.uc.DeleteListBySlug(c, username, listSlug)
	if err != nil {
		helpers.ResponseBuilder(c, http.StatusInternalServerError, "delete list", err.Error(), nil)
		return
	}

	helpers.ResponseBuilder(c, http.StatusOK, "delete list", nil, nil)
}
