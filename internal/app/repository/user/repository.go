package userrepository

import (
	"context"
	"goproject/internal/domain/model"
	"goproject/internal/domain/repository"

	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (repo *UserRepositoryImpl) Create(ctx context.Context, data model.User, tx *gorm.DB) error {
	err := tx.WithContext(ctx).Create(&data).Error
	return err
}

func (repo *UserRepositoryImpl) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	user := new(model.User)

	err := repo.db.WithContext(ctx).First(user, "username = ?", username).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *UserRepositoryImpl) Update(ctx context.Context, data model.User) error {
	err := repo.db.WithContext(ctx).Updates(data).Error
	return err
}
