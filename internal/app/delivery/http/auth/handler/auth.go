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
//
//	@Summary		Create a new account
//	@Description	Create a new account by providing required data. This will automatically create a blog named: "<user's name>'s blog".
//	@Description	User's username and email must be unique. Meaning that there can't be 2 users using the same email/username.
//	@Tags			Auth
//	@Param			Body	body	dto.RegisterRequest	true	"data required to create a new account"
//	@Produce		json
//	@Success		201	{object}	helpers.ResponseWithoutDataAndError
//	@Failure		409	{object}	helpers.ResponseWithError
//	@Failure		400	{object}	helpers.ResponseWithError{error=[]helpers.InputError}
//	@Router			/auth/register [post]
func (handler *UserHandlerImpl) Register(c *gin.Context) {
	var data dto.RegisterRequest

	err := c.ShouldBindJSON(&data)
	if err != nil {
		helpers.ResponseBuilder(c, http.StatusBadRequest, "register", helpers.ValidationError(err), nil)
		return
	}

	ucErr := handler.uc.Register(c, data)
	if ucErr != nil {
		helpers.ResponseBuilder(c, ucErr.Code, "register", ucErr.String(), nil)
		return
	}

	helpers.ResponseBuilder(c, http.StatusCreated, "register", nil, nil)
}

// UserLogin godoc
//
//	@Summary		Login as an existing user
//	@Description	Log in as an existing user by providing a username and password
//
//	@Description	Upon successful login, a JWT will be provided
//
//	@Tags			Auth
//	@Param			Body	body	dto.LoginRequest	true	"data required to login to an existing account"
//	@Produce		json
//	@Success		200	{object}	helpers.ResponseWithData{data=dto.LoginResponse}
//	@Failure		401	{object}	helpers.ResponseWithError
//	@Router			/auth/login [post]
func (handler *UserHandlerImpl) Login(c *gin.Context) {
	var data dto.LoginRequest

	err := c.ShouldBindJSON(&data)
	if err != nil {
		helpers.ResponseBuilder(c, http.StatusBadRequest, "login", helpers.ValidationError(err), nil)
		return
	}

	resp, ucErr := handler.uc.Login(c, data)
	if ucErr != nil {
		helpers.ResponseBuilder(c, ucErr.Code, "login", ucErr.String(), nil)
		return
	}

	helpers.ResponseBuilder(c, http.StatusOK, "login", nil, resp)
}
