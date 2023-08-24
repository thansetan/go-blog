package blogusecase

import (
	"context"
	"errors"
	"fmt"
	"goproject/internal/app/delivery/http/blog/dto"
	"goproject/internal/domain/repository"
	"goproject/internal/helpers"
	"log/slog"
	"net/http"

	"gorm.io/gorm"
)

type BlogUsecase interface {
	UpdateBlogData(ctx context.Context, username string, data dto.UpdateBlogRequest) *helpers.Error
	GetBlogByOwner(ctx context.Context, owner string) (*dto.BlogResponse, *helpers.Error)
}

type blogUsecaseImpl struct {
	repo   repository.BlogRepository
	logger *slog.Logger
}

func NewBlogUsecase(repo repository.BlogRepository, logger *slog.Logger) BlogUsecase {
	return &blogUsecaseImpl{
		repo:   repo,
		logger: logger,
	}
}

func (uc *blogUsecaseImpl) UpdateBlogData(ctx context.Context, username string, data dto.UpdateBlogRequest) *helpers.Error {
	blog, err := uc.repo.FindByOwner(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helpers.ErrorBuilder(http.StatusNotFound, fmt.Sprintf("%s's blog not found", username))
		}
		uc.logger.ErrorContext(ctx, err.Error())
		return helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	blog.Name = data.Name
	blog.Description = data.Description

	err = uc.repo.Update(ctx, *blog)
	if err != nil {
		uc.logger.ErrorContext(ctx, err.Error())
		return helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	return nil
}

func (uc *blogUsecaseImpl) GetBlogByOwner(ctx context.Context, owner string) (*dto.BlogResponse, *helpers.Error) {
	blog, err := uc.repo.FindByOwner(ctx, owner)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.ErrorBuilder(http.StatusNotFound, fmt.Sprintf("%s's blog not found", owner))
		}
		uc.logger.ErrorContext(ctx, err.Error())
		return nil, helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	data := dto.BlogResponse{
		Name:        blog.Name,
		Owner:       blog.Owner,
		Description: blog.Description,
		NumOfPosts:  len(blog.Posts),
	}
	return &data, nil
}
