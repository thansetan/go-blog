package repository

import (
	"context"
	"goproject/internal/domain/model"
)

type PostRepository interface {
	Create(ctx context.Context, data model.Post) error
	FindByBlogID(ctx context.Context, blogID uint) ([]model.Post, error)
	Update(ctx context.Context, data model.Post) error
	FindBySlugAndBlogID(ctx context.Context, slug string, blogID uint) (*model.Post, error)
	Delete(ctx context.Context, data model.Post) error
}
