package userhandler

import (
	"goproject/internal/app/delivery/http/user/dto"
	userusecase "goproject/internal/app/usecase/user"
	"goproject/internal/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	GetMyInformation(c *gin.Context)
	UpdateMyPassword(c *gin.Context)
	UpdateMyInformation(c *gin.Context)
}

type userHandlerImpl struct {
	uc userusecase.UserUsecase
}

func NewUserHandler(usecase userusecase.UserUsecase) UserHandler {
	return &userHandlerImpl{
		uc: usecase,
	}
}

// MyInformation godoc
//
//	@Summary		Get current user information
//	@Description	Get user information about current logged in user
//	@Tags			User
//	@Security		BearerToken
//	@Product		200 {object} map[string]any
//	@Router			/users/me [get]
func (handler *userHandlerImpl) GetMyInformation(c *gin.Context) {
	username := c.GetString("username")

	if username == "" {
		helpers.ResponseBuilder(c, http.StatusUnauthorized, "get my information", "you're not allowed to access this path", nil)
		return
	}

	user, err := handler.uc.GetUserDataByUsername(c, username)
	if err != nil {
		helpers.ResponseBuilder(c, err.Code, "get my information", err.String(), nil)
		return
	}

	helpers.ResponseBuilder(c, http.StatusOK, "get my information", nil, user)
}

// ChangePassword godoc
//
//	@Summary		Change user password
//	@Description	Change user password by providing required data
//	@Tags			User
//	@Param			Body	body	dto.UpdatePasswordRequest	true	"the body to change user's password"
//	@Security		BearerToken
//	@Produce		json
//	@Success		200	{object}	map[string]any
//	@Router			/users/me/update-password [put]
func (handler *userHandlerImpl) UpdateMyPassword(c *gin.Context) {
	var data dto.UpdatePasswordRequest
	username := c.GetString("username")

	if username == "" {
		helpers.ResponseBuilder(c, http.StatusUnauthorized, "update password", "you're not allowed to access this path", nil)
		return
	}

	err := c.ShouldBindJSON(&data)
	if err != nil {
		helpers.ResponseBuilder(c, http.StatusBadRequest, "update password", helpers.ValidationError(err), nil)
		return
	}

	ucErr := handler.uc.ChangePasswordByUsername(c, username, data)
	if ucErr != nil {
		helpers.ResponseBuilder(c, ucErr.Code, "update password", ucErr.String(), nil)
		return
	}

	helpers.ResponseBuilder(c, http.StatusOK, "update password", nil, nil)
}

// UpdateMyInformation godoc
//
//	@Summary		Update current user's information
//	@Description	Update current user's information by providing required data
//	@Tags			User
//	@Param			Body	body	dto.UserUpdateInfoRequest	true	"the body to update user's information"
//	@Security		BearerToken
//	@Produce		json
//	@Success		200	{object}	map[string]any
//	@Router			/users/me [put]
func (handler *userHandlerImpl) UpdateMyInformation(c *gin.Context) {
	var data dto.UserUpdateInfoRequest
	username := c.GetString("username")

	if username == "" {
		helpers.ResponseBuilder(c, http.StatusUnauthorized, "update information", "you're not allowed to access this path", nil)
		return
	}

	err := c.ShouldBindJSON(&data)
	if err != nil {
		helpers.ResponseBuilder(c, http.StatusBadRequest, "update information", helpers.ValidationError(err), nil)
		return
	}

	ucErr := handler.uc.UpdateUserInformation(c, username, data)
	if ucErr != nil {
		helpers.ResponseBuilder(c, ucErr.Code, "update information", ucErr.String(), nil)
		return
	}

	helpers.ResponseBuilder(c, http.StatusOK, "update information", nil, nil)
}
