package userhandler

import (
	"goproject/internal/app/delivery/http/user/dto"
	userusecase "goproject/internal/app/usecase/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	GetMyInformation(c *gin.Context)
	UpdateMyPassword(c *gin.Context)
	UpdateMyInformation(c *gin.Context)
}

type UserHandlerImpl struct {
	uc userusecase.UserUsecase
}

func NewUserHandler(usecase userusecase.UserUsecase) UserHandler {
	return &UserHandlerImpl{
		uc: usecase,
	}
}

// MyInformation godoc
// @Summary Get current user information
// @Description Get user information about current logged in user
// @Tags User
// @Security BearerToken
// @Product 200 {object} map[string]any
// @Router /users/me [get]
func (handler *UserHandlerImpl) GetMyInformation(c *gin.Context) {
	username := c.GetString("username")

	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "username not provided",
		})
		return
	}

	user, err := handler.uc.GetUserDataByUsername(c, username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change user password by providing required data
// @Tags User
// @Param Body body dto.UpdatePasswordRequest true "the body to change user's password"
// @Security BearerToken
// @Produce json
// @Success 200 {object} map[string]any
// @Router /users/me/update-password [put]
func (handler *UserHandlerImpl) UpdateMyPassword(c *gin.Context) {
	var data dto.UpdatePasswordRequest
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

	err = handler.uc.ChangePasswordByUsername(c, username, data)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// UpdateMyInformation godoc
// @Summary Update current user's information
// @Description Update current user's information by providing required data
// @Tags User
// @Param Body body dto.UserUpdateInfoRequest true "the body to update user's information"
// @Security BearerToken
// @Produce json
// @Success 200 {object} map[string]any
// @Router /users/me [put]
func (handler *UserHandlerImpl) UpdateMyInformation(c *gin.Context) {
	var data dto.UserUpdateInfoRequest
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

	err = handler.uc.UpdateUserInformation(c, username, data)
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
