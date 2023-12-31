package postrepository

import (
	"context"
	"goproject/internal/domain/model"
	"goproject/internal/domain/repository"

	"gorm.io/gorm"
)

type postRepositoryImpl struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) repository.PostRepository {
	return &postRepositoryImpl{
		db: db,
	}
}

func (repo *postRepositoryImpl) Create(ctx context.Context, data model.Post) error {
	err := repo.db.WithContext(ctx).Create(&data).Error
	return err
}

func (repo *postRepositoryImpl) FindByBlogID(ctx context.Context, blogID uint) ([]model.Post, error) {
	var posts []model.Post
	err := repo.db.WithContext(ctx).Joins("Blog.User").Find(&posts, "blog_id=?", blogID).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (repo *postRepositoryImpl) Update(ctx context.Context, data model.Post) error {
	err := repo.db.WithContext(ctx).Updates(&data).Error
	return err
}

func (repo *postRepositoryImpl) FindBySlugAndOwner(ctx context.Context, slug, owner string) (*model.Post, error) {
	post := new(model.Post)
	err := repo.db.WithContext(ctx).Joins("Blog.User").Joins("Blog").First(post, "slug=? AND owner=?", slug, owner).Error
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (repo *postRepositoryImpl) Delete(ctx context.Context, post model.Post) error {
	err := repo.db.WithContext(ctx).Delete(&post).Error
	return err
}
