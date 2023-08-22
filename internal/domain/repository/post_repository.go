package repository

import (
	"context"
	"goproject/internal/domain/model"
)

type PostRepository interface {
	Create(ctx context.Context, data model.Post) error
	FindByBlogID(ctx context.Context, blogID uint) ([]model.Post, error)
	Update(ctx context.Context, data model.Post) error
	FindBySlugAndOwner(ctx context.Context, slug, owner string) (*model.Post, error)
	Delete(ctx context.Context, data model.Post) error
}
