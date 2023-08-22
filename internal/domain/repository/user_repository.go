package repository

import (
	"context"
	"goproject/internal/domain/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, data model.User, tx *gorm.DB) error
	Update(ctx context.Context, data model.User) error
	FindByUsername(ctx context.Context, username string) (*model.User, error)
}
