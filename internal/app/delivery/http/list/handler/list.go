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

type listHandlerImpl struct {
	uc listusecase.ListUsecase
}

func NewListHandler(uc listusecase.ListUsecase) ListHandler {
	return &listHandlerImpl{
		uc: uc,
	}
}

//	@CreateNewList	godoc
//	@Summary		Create a new list
//	@Description	Create a new list for current user by providing required data.
//	@Description	Upon successful creation, it will return the newly created list's slug.
//	@Tags			List
//	@Param			Body	body	dto.ListRequest	true	"data required to create a new list"
//	@Security		BearerToken
//	@Produce		json
//	@Success		201	{object}	helpers.ResponseWithData{data=dto.CreateListResponse}
//	@Router			/lists/my [post]
func (handler *listHandlerImpl) CreateNewList(c *gin.Context) {
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

	id, ucErr := handler.uc.CreateNewList(c, data, username)
	if ucErr != nil {
		helpers.ResponseBuilder(c, ucErr.Code, "create list", ucErr.String(), nil)
		return
	}

	helpers.ResponseBuilder(c, http.StatusCreated, "create list", nil, id)
}

//	@AddPostToMyList	godoc
//	@Summary			Add post to current user's list
//	@Description		Add a post to the current user's list by providing the slug of the list to which the user wants to add the post.
//	@Tags				List
//	@Param				username	path	string	true	"blog owner's username"
//	@Param				post_slug	path	string	true	"post's slug"
//	@Param				list_slug	path	string	true	"list slug you want to add this post to"
//	@Security			BearerToken
//	@Produce			json
//	@Success			200	{object}	helpers.ResponseWithoutDataAndError
//	@Failure			404	{object}	helpers.ResponseWithError
//	@Failure			409	{object}	helpers.ResponseWithError
//	@Router				/blog/{username}/posts/{post_slug}/save/{list_slug} [post]
func (handler *listHandlerImpl) AddPostToMyList(c *gin.Context) {
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
		helpers.ResponseBuilder(c, err.Code, "add post to list", err.String(), nil)
		return
	}

	helpers.ResponseBuilder(c, http.StatusOK, "add post to list", nil, nil)
}

//	@GetPostsInMyList	godoc
//	@Summary			Get posts in current user's list
//	@Description		Get posts in my current user's by providing the list's slug.
//	@Tags				List
//	@Param				list_slug	path	string	true	"slug of the list you want to get the post from"
//	@Security			BearerToken
//	@Produce			json
//	@Success			200	{object}	helpers.ResponseWithData{data=helpers.PostsInMyListResponse}
//	@Router				/lists/my/{list_slug} [get]
func (handler *listHandlerImpl) GetPostsInMyListBySlug(c *gin.Context) {
	username := c.GetString("username")
	listSlug := c.Param("list_slug")

	if username == "" {
		helpers.ResponseBuilder(c, http.StatusUnauthorized, "get posts in list", "you're not allowed to access this path", nil)
		return
	}

	list, err := handler.uc.GetPostsInAListBySlug(c, listSlug, username)
	if err != nil {
		helpers.ResponseBuilder(c, err.Code, "get posts in list", err.String(), nil)
		return
	}

	helpers.ResponseBuilder(c, http.StatusOK, "get posts in list", nil, list)
}

//	@GetMyLists		godoc
//	@Summary		Get current user's lists
//	@Description	Get all of current user's lists
//	@Description	Will return an empty array ([]) if the user has no lists.
//	@Tags			List
//	@Security		BearerToken
//	@Produce		json
//	@Success		200	{object}	helpers.ResponseWithData{data=[]helpers.MyListResponse}
//	@Router			/lists/my [get]
func (handler *listHandlerImpl) GetMyLists(c *gin.Context) {
	username := c.GetString("username")

	if username == "" {
		helpers.ResponseBuilder(c, http.StatusUnauthorized, "get lists", "you're not allowed to access this path", nil)
		return
	}

	lists, err := handler.uc.GetMyLists(c, username)
	if err != nil {
		helpers.ResponseBuilder(c, err.Code, "get lists", err.String(), nil)
		return
	}

	helpers.ResponseBuilder(c, http.StatusOK, "get lists", nil, lists)
}

//	@UpdateMyListInformation	godoc
//	@Summary					Update/modify current user's list information
//	@Description				Update/modify current user's list information (name and description) by providing the list's slug.
//	@Tags						List
//	@Param						list_slug	path	string			true	"list slug you want to edit"
//	@Param						body		body	dto.ListRequest	strue	"body required to update/modify list information"
//	@Security					BearerToken
//	@Produce					json
//	@Success					200	{object}	helpers.ResponseWithoutDataAndError
//	@Failure					404	{object}	helpers.ResponseWithError
//	@Router						/lists/my/{list_slug} [PUT]
func (handler *listHandlerImpl) UpdateMyListInformationBySlug(c *gin.Context) {
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

	ucErr := handler.uc.UpdateListInformation(c, data, username, listSlug)
	if ucErr != nil {
		helpers.ResponseBuilder(c, ucErr.Code, "update list", ucErr.String(), nil)
		return
	}

	helpers.ResponseBuilder(c, http.StatusOK, "update list", nil, nil)
}

//	@RemovePostFromUserList	godoc
//	@Summary				Remove a post from current user's list
//	@Description			Remove a post from current user's list by providing the list slug and post slug you want to remove.
//	@Tags					List
//	@Param					list_slug	path	string	true	"list slug you want to remove post from"
//	@Param					post_slug	path	string	true	"post slug you want to remove"
//	@Security				BearerToken
//	@Produce				json
//	@Success				200	{object}	helpers.ResponseWithoutDataAndError
//	@Failure				404	{object}	helpers.ResponseWithError
//	@Router					/lists/my/{list_slug}/{post_slug} [DELETE]
func (handler *listHandlerImpl) RemovePostFromMyList(c *gin.Context) {
	username := c.GetString("username")
	listSlug := c.Param("list_slug")
	postSlug := c.Param("post_slug")

	if username == "" {
		helpers.ResponseBuilder(c, http.StatusUnauthorized, "remove post from list", "you're not allowed to access this path", nil)
		return
	}

	err := handler.uc.RemovePostFromList(c, username, postSlug, listSlug)
	if err != nil {
		helpers.ResponseBuilder(c, err.Code, "remove post from list", err.String(), nil)
		return
	}

	helpers.ResponseBuilder(c, http.StatusOK, "remove post from list", nil, nil)
}

//	@DeleteMyList	godoc
//	@Summary		Delete current user's list
//	@Description	Delete current user's list by providing the list slug
//	@Tags			List
//	@Param			list_slug	path	string	true	"slug of the list you want to remove"
//	@Security		BearerToken
//	@Produce		json
//	@Success		200	{object}	helpers.ResponseWithoutDataAndError
//	@Failure		404	{object}	helpers.ResponseWithError
//	@Router			/lists/my/{list_slug} [DELETE]
func (handler *listHandlerImpl) DeleteMyListBySlug(c *gin.Context) {
	username := c.GetString("username")
	listSlug := c.Param("list_slug")

	if username == "" {
		helpers.ResponseBuilder(c, http.StatusUnauthorized, "delete list", "you're not allowed to access this path", nil)
		return
	}

	err := handler.uc.DeleteListBySlug(c, username, listSlug)
	if err != nil {
		helpers.ResponseBuilder(c, err.Code, "delete list", err.String(), nil)
		return
	}

	helpers.ResponseBuilder(c, http.StatusOK, "delete list", nil, nil)
}
