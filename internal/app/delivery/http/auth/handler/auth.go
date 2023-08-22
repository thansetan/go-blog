package authhandler

import (
	"goproject/internal/app/delivery/http/auth/dto"
	authusecase "goproject/internal/app/usecase/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
}

type UserHandlerImpl struct {
	uc authusecase.AuthUsecase
}

func NewAuthHandler(usecase authusecase.AuthUsecase) UserHandler {
	return &UserHandlerImpl{
		uc: usecase,
	}
}

// UserRegister godoc
// @Summary Create a user account
// @Description Create a new account by providing required data
// @Tags Auth
// @Param Body body dto.RegisterRequest true "the body to register a user"
// @Produce json
// @Success 201 {object} map[string]any
// @Router /auth/register [post]
func (handler *UserHandlerImpl) Register(c *gin.Context) {
	var data dto.RegisterRequest

	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
		})
		return
	}

	err = handler.uc.Register(c, data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"err":     err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
	})
}

// UserLogin godoc
// @Summary Login as a user
// @Description Logging in by providing required data to get JWT
// @Tags Auth
// @Param Body body dto.LoginRequest true "the body to login as a user"
// @Produce json
// @Success 200 {object} map[string]any
// @Router /auth/login [post]
func (handler *UserHandlerImpl) Login(c *gin.Context) {
	var data dto.LoginRequest

	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
		})
		return
	}

	resp, err := handler.uc.Login(c, data)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
