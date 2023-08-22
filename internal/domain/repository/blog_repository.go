package repository

import (
	"context"
	"goproject/internal/domain/model"

	"gorm.io/gorm"
)

type BlogRepository interface {
	Create(ctx context.Context, data model.Blog, tx *gorm.DB) error
	Update(ctx context.Context, data model.Blog) error
	FindByOwner(ctx context.Context, username string) (*model.Blog, error)
}
