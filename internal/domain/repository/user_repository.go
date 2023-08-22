package repository

import (
	"context"
	"goproject/internal/domain/model"
)

type UserRepository interface {
	Create(ctx context.Context, data model.User) error
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	UpdateByUsername(ctx context.Context, data model.User) error
}
