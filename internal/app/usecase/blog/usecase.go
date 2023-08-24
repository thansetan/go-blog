package blogusecase

import (
	"context"
	"goproject/internal/app/delivery/http/blog/dto"
	"goproject/internal/domain/repository"
)

type BlogUsecase interface {
	UpdateBlogData(ctx context.Context, username string, data dto.UpdateBlogRequest) error
	GetBlogByOwner(ctx context.Context, owner string) (*dto.BlogResponse, error)
}

type BlogUsecaseImpl struct {
	repo repository.BlogRepository
}

func NewBlogUsecase(repo repository.BlogRepository) BlogUsecase {
	return &BlogUsecaseImpl{
		repo: repo,
	}
}

func (uc *BlogUsecaseImpl) UpdateBlogData(ctx context.Context, username string, data dto.UpdateBlogRequest) error {
	blog, err := uc.repo.FindByOwner(ctx, username)
	if err != nil {
		return err
	}

	blog.Name = data.Name
	blog.Description = data.Description

	err = uc.repo.Update(ctx, *blog)
	if err != nil {
		return err
	}

	return nil
}

func (uc *BlogUsecaseImpl) GetBlogByOwner(ctx context.Context, owner string) (*dto.BlogResponse, error) {
	blog, err := uc.repo.FindByOwner(ctx, owner)
	if err != nil {
		return nil, err
	}

	data := dto.BlogResponse{
		Name:        blog.Name,
		Owner:       blog.Owner,
		Description: blog.Description,
		NumOfPosts:  len(blog.Posts),
	}
	return &data, nil
}
