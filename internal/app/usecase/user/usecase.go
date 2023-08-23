package userusecase

import (
	"context"
	"goproject/internal/app/delivery/http/user/dto"
	"goproject/internal/domain/repository"
	"goproject/internal/utils"
)

type UserUsecase interface {
	GetUserDataByUsername(ctx context.Context, username string) (*dto.UserResponse, error)
	ChangePasswordByUsername(ctx context.Context, username string, data dto.UpdatePasswordRequest) error
	UpdateUserInformation(ctx context.Context, username string, data dto.UserUpdateInfoRequest) error
}

type UserUsecaseImpl struct {
	repo repository.UserRepository
}

func NewUserUsecase(repository repository.UserRepository) UserUsecase {
	return &UserUsecaseImpl{
		repo: repository,
	}
}

func (uc *UserUsecaseImpl) GetUserDataByUsername(ctx context.Context, username string) (*dto.UserResponse, error) {
	data, err := uc.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	user := &dto.UserResponse{
		Name:     data.Name,
		Username: data.Username,
		Email:    data.Email,
	}

	return user, nil
}

func (uc *UserUsecaseImpl) ChangePasswordByUsername(ctx context.Context, username string, data dto.UpdatePasswordRequest) error {
	user, err := uc.repo.FindByUsername(ctx, username)
	if err != nil {
		return err
	}

	err = utils.IsValidPassword(user.Password, data.OldPassword)
	if err != nil {
		return err
	}

	newPassword, err := utils.HashPassword(data.NewPassword)
	if err != nil {
		return err
	}

	user.Password = newPassword

	err = uc.repo.Update(ctx, *user)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UserUsecaseImpl) UpdateUserInformation(ctx context.Context, username string, data dto.UserUpdateInfoRequest) error {
	user, err := uc.repo.FindByUsername(ctx, username)
	if err != nil {
		return err
	}

	user.Name = data.Name
	user.Email = data.Email

	err = uc.repo.Update(ctx, *user)
	if err != nil {
		return err
	}

	return nil
}
