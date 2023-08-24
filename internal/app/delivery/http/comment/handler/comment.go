package commenthandler

import (
	"goproject/internal/app/delivery/http/comment/dto"
	commentusecase "goproject/internal/app/usecase/comment"
	"goproject/internal/helpers"
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

type commentHandlerImpl struct {
	uc commentusecase.CommentUsecase
}

func NewCommentHandler(uc commentusecase.CommentUsecase) CommentHandler {
	return &commentHandlerImpl{
		uc: uc,
	}
}

//	@CreateComment	godoc
//	@Summary		Create a comment on a blog post
//	@Description	Create a new comment on a blog post.
//	@Description	Upon successful creation, it will returns the newly created comment's ID.
//	@Tags			Comment
//	@Param			Body		body	dto.CommentRequest	true	"data required to create a comment"
//	@Param			username	path	string				true	"blog owner's username"
//	@Param			post_slug	path	string				true	"post's slug"
//	@Security		BearerToken
//	@Produce		json
//	@Success		201	{object}	helpers.ResponseWithData{data=dto.CreateCommentResponse}
//	@Failure		404	{object}	helpers.ResponseWithError
//	@Router			/blog/{username}/posts/{post_slug}/comments [post]
func (handler *commentHandlerImpl) CreateComment(c *gin.Context) {
	var data dto.CommentRequest
	blogOwner := c.Param("username")
	postSlug := c.Param("post_slug")
	username := c.GetString("username")

	if username == "" {
		helpers.ResponseBuilder(c, http.StatusUnauthorized, "create comment", "you're not allowed to access this path", nil)
		return
	}

	err := c.ShouldBindJSON(&data)
	if err != nil {
		helpers.ResponseBuilder(c, http.StatusBadRequest, "create comment", helpers.ValidationError(err), nil)
		return
	}

	commentID, ucErr := handler.uc.CreateComment(c, data, username, blogOwner, postSlug)
	if ucErr != nil {
		helpers.ResponseBuilder(c, ucErr.Code, "create comment", ucErr.String(), nil)
		return
	}

	helpers.ResponseBuilder(c, http.StatusCreated, "create comment", nil, commentID)
}

//	@GetMyComments	Godoc
//	@Summary		Get current user's comments
//	@Description	Get current user's comments on all posts.
//	@Tags			Comment
//	@Security		BearerToken
//	@Produce		json
//	@Success		200	{object}	helpers.ResponseWithData{data=dto.CommentResponse}
//	@Router			/my/comments [get]
func (handler *commentHandlerImpl) GetMyComments(c *gin.Context) {
	username := c.GetString("username")
	if username == "" {
		helpers.ResponseBuilder(c, http.StatusUnauthorized, "get comments", "you're not allowed to access this path", nil)
		return
	}

	comments, err := handler.uc.GetCommentsByUsername(c, username)
	if err != nil {
		helpers.ResponseBuilder(c, err.Code, "get comments", err.String(), nil)
		return
	}

	helpers.ResponseBuilder(c, http.StatusOK, "get comments", nil, comments)
}

//	@GetCommentOnAPost	Godoc
//	@Summary			Get comments on a post
//	@Description		Get all comments on a post by post's URL.
//	@Tags				Comment
//	@Param				username	path	string	true	"blog owner's username"
//	@Param				post_slug	path	string	true	"post's slug"
//	@Produce			json
//	@Success			200	{object}	helpers.ResponseWithData{data=[]dto.CommentResponse}
//	@Router				/blog/{username}/posts/{post_slug}/comments [get]
func (handler *commentHandlerImpl) GetCommentsOnAPost(c *gin.Context) {
	blogOwner := c.Param("username")
	postSlug := c.Param("post_slug")

	comments, err := handler.uc.GetCommentsByBlogOwnerAndPostSlug(c, blogOwner, postSlug)
	if err != nil {
		helpers.ResponseBuilder(c, err.Code, "get comments", err.String(), nil)
		return
	}

	helpers.ResponseBuilder(c, http.StatusOK, "get comments", nil, comments)
}

//	@DeleteCommentByID	Godoc
//	@Summary			Delete a comment
//	@Description		Delete a comment on a post by comment's ID.
//	@Description		A non-blog-owner user can only delete their own comment.
//	@Description		Blog's owner is allowed to delete ANY comment on their posts.
//	@Tags				Comment
//	@Param				username	path	string	true	"blog owner's username"
//	@Param				post_slug	path	string	true	"post's slug"
//	@Param				comment_id	path	int		true	"comment's ID"
//	@Security			BearerToken
//	@Produce			json
//	@Success			200	{object}	helpers.ResponseWithoutDataAndError
//	@Failure			404	{object}	helpers.ResponseWithError
//	@Router				/blog/{username}/posts/{post_slug}/comments/{comment_id} [delete]
func (handler *commentHandlerImpl) DeleteCommentByID(c *gin.Context) {
	blogOwner := c.Param("username")
	postSlug := c.Param("post_slug")
	commentID := c.Param("comment_id")
	username := c.GetString("username")

	if username == "" {
		helpers.ResponseBuilder(c, http.StatusUnauthorized, "delete comment", "you're not allowed to access this path", nil)
		return
	}

	err := handler.uc.DeleteCommentOnAPosst(c, username, blogOwner, postSlug, commentID)
	if err != nil {
		helpers.ResponseBuilder(c, err.Code, "delete comment", err.String(), nil)
		return
	}

	helpers.ResponseBuilder(c, http.StatusOK, "delete comment", nil, nil)
}

//	@EditMyCommentOnAPost	Godoc
//	@Summary				Edit current user's comment on a post
//	@Description			Edit current user's comment on a post by providing comment ID.
//	@Tags					Comment
//	@Param					username	path	string				true	"blog owner's username"
//	@Param					post_slug	path	string				true	"post's slug"
//	@Param					Body		body	dto.CommentRequest	true	"data required to modify comment"
//	@Param					comment_id	path	int					true	"comment's ID"
//	@Security				BearerToken
//	@Produce				json
//	@Success				200	{object}	helpers.ResponseWithoutDataAndError
//	@Failure				404	{object}	helpers.ResponseWithError
//	@Router					/blog/{username}/posts/{post_slug}/comments/{comment_id} [put]
func (handler *commentHandlerImpl) EditMyCommentOnAPost(c *gin.Context) {
	var data dto.CommentRequest
	blogOwner := c.Param("username")
	postSlug := c.Param("post_slug")
	commentID := c.Param("comment_id")
	username := c.GetString("username")

	if username == "" {
		helpers.ResponseBuilder(c, http.StatusUnauthorized, "edit comment", "you're not allowed to access this path", nil)
		return
	}

	err := c.ShouldBindJSON(&data)
	if err != nil {
		helpers.ResponseBuilder(c, http.StatusBadRequest, "edit comment", helpers.ValidationError(err), nil)
		return
	}

	ucErr := handler.uc.UpdateCommentOnAPost(c, username, blogOwner, postSlug, commentID, data)
	if ucErr != nil {
		helpers.ResponseBuilder(c, ucErr.Code, "edit comment", ucErr.String(), nil)
		return
	}

	helpers.ResponseBuilder(c, http.StatusOK, "edit comment", nil, nil)
}
