package commentrepository

import (
	"context"
	"goproject/internal/domain/model"
	"goproject/internal/domain/repository"

	"gorm.io/gorm"
)

type CommentRepositoryImpl struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) repository.CommentRepository {
	return &CommentRepositoryImpl{
		db: db,
	}
}

func (repo *CommentRepositoryImpl) Create(ctx context.Context, data model.Comment) error {
	err := repo.db.WithContext(ctx).Create(&data).Error
	return err
}

func (repo *CommentRepositoryImpl) FindCommentByUsername(ctx context.Context, username string) ([]model.Comment, error) {
	var comments []model.Comment

	err := repo.db.WithContext(ctx).Joins("Post.Blog").Find(&comments, "commenter=?", username).Error
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (repo *CommentRepositoryImpl) FindCommentByPostID(ctx context.Context, PostID uint) ([]model.Comment, error) {
	var comments []model.Comment

	err := repo.db.WithContext(ctx).Find(&comments, "post_id=?", PostID).Error
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (repo *CommentRepositoryImpl) Delete(ctx context.Context, data model.Comment) error {
	err := repo.db.WithContext(ctx).Delete(&data).Error
	return err
}

func (repo *CommentRepositoryImpl) FindCommentByID(ctx context.Context, commentID uint) (*model.Comment, error) {
	comment := new(model.Comment)
	err := repo.db.WithContext(ctx).First(comment, "id=?", commentID).Error
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (repo *CommentRepositoryImpl) Update(ctx context.Context, data model.Comment) error {
	err := repo.db.WithContext(ctx).Updates(&data).Error
	return err
}
