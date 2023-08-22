package blogrepository

import (
	"context"
	"goproject/internal/domain/model"
	"goproject/internal/domain/repository"

	"gorm.io/gorm"
)

type BlogRepositoryImpl struct {
	db *gorm.DB
}

func NewBlogRepository(db *gorm.DB) repository.BlogRepository {
	return &BlogRepositoryImpl{
		db: db,
	}
}

func (repo *BlogRepositoryImpl) Create(ctx context.Context, data model.Blog, tx *gorm.DB) error {
	err := tx.WithContext(ctx).Create(&data).Error
	return err
}

func (repo *BlogRepositoryImpl) FindByOwner(ctx context.Context, owner string) (*model.Blog, error) {
	blog := new(model.Blog)

	err := repo.db.WithContext(ctx).First(blog, "owner = ?", owner).Error
	if err != nil {
		return nil, err
	}

	return blog, nil
}

func (repo *BlogRepositoryImpl) Update(ctx context.Context, data model.Blog) error {
	err := repo.db.Updates(&data).Error
	return err
}