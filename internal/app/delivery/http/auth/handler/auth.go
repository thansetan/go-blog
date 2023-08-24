package authhandler

import (
	"goproject/internal/app/delivery/http/auth/dto"
	authusecase "goproject/internal/app/usecase/auth"
	"goproject/internal/helpers"
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
// @Description Create a new account by providing required data. This will automatically create a blog named: <name>'s Blog
// @Tags Auth
// @Param Body body dto.RegisterRequest true "the body to register a user"
// @Produce json
// @Success 201 {object} map[string]any
// @Router /auth/register [post]
func (handler *UserHandlerImpl) Register(c *gin.Context) {
	var data dto.RegisterRequest

	err := c.ShouldBindJSON(&data)
	if err != nil {
		helpers.ResponseBuilder(c, http.StatusBadRequest, "register", helpers.ValidationError(err), nil)
		return
	}

	err = handler.uc.Register(c, data)
	if err != nil {
		helpers.ResponseBuilder(c, http.StatusInternalServerError, "register", err.Error(), nil)
		return
	}

	helpers.ResponseBuilder(c, http.StatusCreated, "register", nil, nil)
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
		helpers.ResponseBuilder(c, http.StatusBadRequest, "login", helpers.ValidationError(err), nil)
		return
	}

	resp, err := handler.uc.Login(c, data)
	if err != nil {
		helpers.ResponseBuilder(c, http.StatusInternalServerError, "login", err.Error(), nil)
		return
	}

	helpers.ResponseBuilder(c, http.StatusOK, "login", nil, resp)
}
