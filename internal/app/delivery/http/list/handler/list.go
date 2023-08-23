package listhandler

import (
	"goproject/internal/app/delivery/http/list/dto"
	listusecase "goproject/internal/app/usecase/list"
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
// @Success 201 {objects} map[string]any
// @Router /lists/my [post]
func (handler *ListHandlerImpl) CreateNewList(c *gin.Context) {
	var data dto.ListRequest
	username := c.GetString("username")

	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "username not provided",
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

	err = handler.uc.CreateNewList(c, data, username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
	})
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
// @Success 200 {objects} map[string]any
// @Router /blog/{username}/posts/{post_slug}/save/{list_slug} [post]
func (handler *ListHandlerImpl) AddPostToMyList(c *gin.Context) {
	username := c.GetString("username")
	postSlug := c.Param("post_slug")
	blogOwner := c.Param("username")
	listSlug := c.Param("list_slug")

	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "username not provided",
		})
		return
	}

	err := handler.uc.AddPostToMyList(c, listSlug, username, blogOwner, postSlug)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
	})
}

// @GetPostsInMyList godoc
// @Summary Get posts in my list
// @Description Get posts in my list by providing required data
// @Tags List
// @Param list_slug path string true "list slug you want to get"
// @Security BearerToken
// @Produce json
// @Success 200 {objects} map[string]any
// @Router /lists/my/{list_slug} [get]
func (handler *ListHandlerImpl) GetPostsInMyListBySlug(c *gin.Context) {
	username := c.GetString("username")
	listSlug := c.Param("list_slug")

	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "username not provided",
		})
		return
	}

	list, err := handler.uc.GetPostsInAListBySlug(c, listSlug, username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    list,
	})
}

// @GetMyLists godoc
// @Summary Get current user's lists
// @Description Get current user's lists by providing required data
// @Tags List
// @Security BearerToken
// @Produce json
// @Success 200 {objects} map[string]any
// @Router /lists/my [get]
func (handler *ListHandlerImpl) GetMyLists(c *gin.Context) {
	username := c.GetString("username")

	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "username not provided",
		})
		return
	}

	lists, err := handler.uc.GetMyLists(c, username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    lists,
	})
}

// @UpdateMyListInformation godoc
// @Summary Update current user's list by ID
// @Description Update current user's list information (name and description) by providing the list ID
// @Tags List
// @Param list_slug path string true "list slug you want to edit"
// @Param body body dto.ListRequest strue "body to update"
// @Security BearerToken
// @Produce json
// @Success 200 {objects} map[string]any
// @Router /lists/my/{list_slug} [PUT]
func (handler *ListHandlerImpl) UpdateMyListInformationBySlug(c *gin.Context) {
	var data dto.ListRequest
	username := c.GetString("username")
	listSlug := c.Param("list_slug")

	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "username not provided",
		})
		return
	}

	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	err = handler.uc.UpdateListInformation(c, data, username, listSlug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// @RemovePostFromUserList godoc
// @Summary Remove post from current user's list
// @Description Remove post from current user's list by providing the list ID and post slug
// @Tags List
// @Param list_slug path string true "list slug you want to remove post from"
// @Param post_slug path string true "post slug you want to delete"
// @Security BearerToken
// @Produce json
// @Success 200 {objects} map[string]any
// @Router /lists/my/{list_slug}/{post_slug} [DELETE]
func (handler *ListHandlerImpl) RemovePostFromMyList(c *gin.Context) {
	username := c.GetString("username")
	listSlug := c.Param("list_slug")
	postSlug := c.Param("post_slug")

	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "username not provided",
		})
		return
	}

	err := handler.uc.RemovePostFromList(c, username, postSlug, listSlug)
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

// @DeleteMyList godoc
// @Summary Delete current user's list by ID
// @Description Delete current user's list by providing the list ID
// @Tags List
// @Param list_slug path string true "list slug you want to remove"
// @Security BearerToken
// @Produce json
// @Success 200 {objects} map[string]any
// @Router /lists/my/{list_slug} [DELETE]
func (handler *ListHandlerImpl) DeleteMyListBySlug(c *gin.Context) {
	username := c.GetString("username")
	listSlug := c.Param("list_slug")

	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "username not provided",
		})
		return
	}

	err := handler.uc.DeleteListBySlug(c, username, listSlug)
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
