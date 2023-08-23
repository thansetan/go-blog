package authusecase

import (
	"context"
	"fmt"
	"goproject/internal/app/delivery/http/auth/dto"
	"goproject/internal/domain/model"
	"goproject/internal/domain/repository"
	"goproject/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthUsecase interface {
	Register(ctx context.Context, data dto.RegisterRequest) error
	Login(ctx context.Context, data dto.LoginRequest) (*dto.LoginResponse, error)
}

type AuthUsecaseImpl struct {
	userRepo repository.UserRepository
	blogRepo repository.BlogRepository
	db       *gorm.DB
}

func NewAuthUsecase(userRepo repository.UserRepository, blogRepo repository.BlogRepository, db *gorm.DB) AuthUsecase {
	return &AuthUsecaseImpl{
		userRepo: userRepo,
		blogRepo: blogRepo,
		db:       db,
	}
}

func (uc *AuthUsecaseImpl) Register(ctx context.Context, data dto.RegisterRequest) error {
	password, err := utils.HashPassword(data.Password)
	if err != nil {
		return err
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

	// fmt.Println(blogData, userData)
	tx := uc.db.Begin()

	err = uc.userRepo.Create(ctx, userData, tx)
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return err
	}

	err = uc.blogRepo.Create(ctx, blogData, tx)
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
	}

	tx.Commit()

	return nil
}

func (uc *AuthUsecaseImpl) Login(ctx context.Context, data dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := uc.userRepo.FindByUsername(ctx, data.Username)
	if err != nil {
		return nil, err
	}

	err = utils.IsValidPassword(user.Password, data.Password)
	if err != nil {
		return nil, err
	}

	resp := new(dto.LoginResponse)

	claims := jwt.MapClaims{
		"username": user.Username,
	}

	token, err := utils.GenerateJWT(claims)
	if err != nil {
		return nil, err
	}

	resp.Token = token

	return resp, nil
}
