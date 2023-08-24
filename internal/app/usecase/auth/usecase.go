package authusecase

import (
	"context"
	"errors"
	"fmt"
	"goproject/internal/app/delivery/http/auth/dto"
	"goproject/internal/domain/model"
	"goproject/internal/domain/repository"
	"goproject/internal/helpers"
	"goproject/internal/utils"
	"log/slog"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthUsecase interface {
	Register(ctx context.Context, data dto.RegisterRequest) *helpers.Error
	Login(ctx context.Context, data dto.LoginRequest) (*dto.LoginResponse, *helpers.Error)
}

type AuthUsecaseImpl struct {
	userRepo repository.UserRepository
	blogRepo repository.BlogRepository
	db       *gorm.DB
	logger   *slog.Logger
}

func NewAuthUsecase(userRepo repository.UserRepository, blogRepo repository.BlogRepository, db *gorm.DB, logger *slog.Logger) AuthUsecase {
	return &AuthUsecaseImpl{
		userRepo: userRepo,
		blogRepo: blogRepo,
		db:       db,
		logger:   logger,
	}
}

func (uc *AuthUsecaseImpl) Register(ctx context.Context, data dto.RegisterRequest) *helpers.Error {
	password, err := utils.HashPassword(data.Password)
	if err != nil {
		return helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	userData := model.User{
		Email:    data.Email,
		Password: password,
		Name:     data.Name,
		Username: data.Username,
	}

	blogData := model.Blog{
		Name:        fmt.Sprintf("%s's Blog", data.Name),
		Description: fmt.Sprintf("%s's blog description", data.Name),
		Owner:       data.Username,
	}

	tx := uc.db.Begin()

	err = uc.userRepo.Create(ctx, userData, tx)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return helpers.ErrorBuilder(http.StatusConflict, "username/email already used")
		}
		uc.logger.ErrorContext(ctx, err.Error())
		return helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	err = uc.blogRepo.Create(ctx, blogData, tx)
	if err != nil {
		tx.Rollback()
		uc.logger.ErrorContext(ctx, err.Error())
		return helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	tx.Commit()

	return nil
}

func (uc *AuthUsecaseImpl) Login(ctx context.Context, data dto.LoginRequest) (*dto.LoginResponse, *helpers.Error) {
	user, err := uc.userRepo.FindByUsername(ctx, data.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.ErrorBuilder(http.StatusUnauthorized, "invalid username/password")
		}
		uc.logger.ErrorContext(ctx, err.Error())
		return nil, helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	err = utils.IsValidPassword(user.Password, data.Password)
	if err != nil {
		return nil, helpers.ErrorBuilder(http.StatusUnauthorized, "invalid username/password")
	}

	resp := new(dto.LoginResponse)

	claims := jwt.MapClaims{
		"username": user.Username,
	}

	token, err := utils.GenerateJWT(claims)
	if err != nil {
		uc.logger.ErrorContext(ctx, err.Error())
		return nil, helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	resp.Token = token

	return resp, nil
}
