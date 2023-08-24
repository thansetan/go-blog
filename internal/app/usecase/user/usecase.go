package userusecase

import (
	"context"
	"errors"
	"fmt"
	"goproject/internal/app/delivery/http/user/dto"
	"goproject/internal/domain/repository"
	"goproject/internal/helpers"
	"goproject/internal/utils"
	"log/slog"
	"net/http"

	"gorm.io/gorm"
)

type UserUsecase interface {
	GetUserDataByUsername(ctx context.Context, username string) (*dto.UserResponse, *helpers.Error)
	ChangePasswordByUsername(ctx context.Context, username string, data dto.UpdatePasswordRequest) *helpers.Error
	UpdateUserInformation(ctx context.Context, username string, data dto.UserUpdateInfoRequest) *helpers.Error
}

type userUsecaseImpl struct {
	repo   repository.UserRepository
	logger *slog.Logger
}

func NewUserUsecase(repository repository.UserRepository, logger *slog.Logger) UserUsecase {
	return &userUsecaseImpl{
		repo:   repository,
		logger: logger,
	}
}

func (uc *userUsecaseImpl) GetUserDataByUsername(ctx context.Context, username string) (*dto.UserResponse, *helpers.Error) {
	data, err := uc.repo.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.ErrorBuilder(http.StatusNotFound, fmt.Sprintf("%s not found", username))
		}
		uc.logger.ErrorContext(ctx, err.Error())
		return nil, helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	user := &dto.UserResponse{
		Name:     data.Name,
		Username: data.Username,
		Email:    data.Email,
	}

	return user, nil
}

func (uc *userUsecaseImpl) ChangePasswordByUsername(ctx context.Context, username string, data dto.UpdatePasswordRequest) *helpers.Error {
	user, err := uc.repo.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helpers.ErrorBuilder(http.StatusNotFound, fmt.Sprintf("%s not found", username))
		}
		uc.logger.ErrorContext(ctx, err.Error())
		return helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	err = utils.IsValidPassword(user.Password, data.OldPassword)
	if err != nil {
		return helpers.ErrorBuilder(http.StatusUnauthorized, "old password you provided is incorrect")
	}

	newPassword, err := utils.HashPassword(data.NewPassword)
	if err != nil {
		return helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	user.Password = newPassword

	err = uc.repo.Update(ctx, *user)
	if err != nil {
		uc.logger.ErrorContext(ctx, err.Error())
		return helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}
	return nil
}

func (uc *userUsecaseImpl) UpdateUserInformation(ctx context.Context, username string, data dto.UserUpdateInfoRequest) *helpers.Error {
	user, err := uc.repo.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helpers.ErrorBuilder(http.StatusNotFound, fmt.Sprintf("%s not found", username))
		}
		uc.logger.ErrorContext(ctx, err.Error())
		return helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	user.Name = data.Name
	user.Email = data.Email

	err = uc.repo.Update(ctx, *user)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return helpers.ErrorBuilder(http.StatusConflict, "email already used")
		}
		uc.logger.ErrorContext(ctx, err.Error())
		return helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	return nil
}
