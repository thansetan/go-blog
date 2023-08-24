package blogrepository

import (
	"context"
	"goproject/internal/domain/model"
	"goproject/internal/domain/repository"

	"gorm.io/gorm"
)

type blogRepositoryImpl struct {
	db *gorm.DB
}

func NewBlogRepository(db *gorm.DB) repository.BlogRepository {
	return &blogRepositoryImpl{
		db: db,
	}
}

func (repo *blogRepositoryImpl) Create(ctx context.Context, data model.Blog, tx *gorm.DB) error {
	err := tx.WithContext(ctx).Create(&data).Error
	return err
}

func (repo *blogRepositoryImpl) FindByOwner(ctx context.Context, owner string) (*model.Blog, error) {
	blog := new(model.Blog)

	err := repo.db.WithContext(ctx).Preload("User").Preload("Posts").First(blog, "owner = ?", owner).Error
	if err != nil {
		return nil, err
	}

	return blog, nil
}

func (repo *blogRepositoryImpl) Update(ctx context.Context, data model.Blog) error {
	newData := map[string]any{
		"name":        data.Name,
		"description": data.Description,
	}

	err := repo.db.WithContext(ctx).Model(&data).Updates(newData).Error
	return err
}
