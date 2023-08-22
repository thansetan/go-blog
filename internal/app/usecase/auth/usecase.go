package authusecase

import (
	"context"
	"fmt"
	"goproject/internal/app/delivery/http/auth/dto"
	"goproject/internal/domain/model"
	"goproject/internal/domain/repository"
	"goproject/internal/utils"

	"github.com/golang-jwt/jwt/v5"
)

type AuthUsecase interface {
	Register(ctx context.Context, data dto.RegisterRequest) error
	Login(ctx context.Context, data dto.LoginRequest) (*dto.LoginResponse, error)
}

type AuthUsecaseImpl struct {
	repo repository.UserRepository
}

func NewAuthUsecase(repo repository.UserRepository) AuthUsecase {
	return &AuthUsecaseImpl{
		repo: repo,
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
	err = uc.repo.Create(ctx, userData)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (uc *AuthUsecaseImpl) Login(ctx context.Context, data dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := uc.repo.FindByUsername(ctx, data.Username)
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
