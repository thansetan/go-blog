package commentrepository

import (
	"context"
	"goproject/internal/domain/model"
	"goproject/internal/domain/repository"

	"gorm.io/gorm"
)

type commentRepositoryImpl struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) repository.CommentRepository {
	return &commentRepositoryImpl{
		db: db,
	}
}

func (repo *commentRepositoryImpl) Create(ctx context.Context, data model.Comment) (uint, error) {
	err := repo.db.WithContext(ctx).Create(&data).Error
	return data.ID, err
}

func (repo *commentRepositoryImpl) FindCommentByUsername(ctx context.Context, username string) ([]model.Comment, error) {
	var comments []model.Comment

	err := repo.db.WithContext(ctx).Joins("Post.Blog").Find(&comments, "commenter=?", username).Error
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (repo *commentRepositoryImpl) FindCommentByPostID(ctx context.Context, PostID uint) ([]model.Comment, error) {
	var comments []model.Comment

	err := repo.db.WithContext(ctx).Find(&comments, "post_id=?", PostID).Error
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (repo *commentRepositoryImpl) Delete(ctx context.Context, data model.Comment) error {
	err := repo.db.WithContext(ctx).Delete(&data).Error
	return err
}

func (repo *commentRepositoryImpl) FindCommentByID(ctx context.Context, commentID uint) (*model.Comment, error) {
	comment := new(model.Comment)
	err := repo.db.WithContext(ctx).First(comment, "id=?", commentID).Error
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (repo *commentRepositoryImpl) Update(ctx context.Context, data model.Comment) error {
	err := repo.db.WithContext(ctx).Updates(&data).Error
	return err
}
