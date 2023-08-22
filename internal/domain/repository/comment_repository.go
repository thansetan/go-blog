package repository

import (
	"context"
	"goproject/internal/domain/model"
)

type CommentRepository interface {
	Create(ctx context.Context, data model.Comment) error
	FindCommentByUsername(ctx context.Context, username string) ([]model.Comment, error)
	FindCommentByPostID(ctx context.Context, PostID uint) ([]model.Comment, error)
	Delete(ctx context.Context, data model.Comment) error
	FindCommentByID(ctx context.Context, commentID uint) (*model.Comment, error)
	Update(ctx context.Context, data model.Comment) error
}
