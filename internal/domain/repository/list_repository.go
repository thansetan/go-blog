package repository

import (
	"context"
	"goproject/internal/domain/model"
)

type ListRepository interface {
	Create(ctx context.Context, data model.List) error
	AddPostToList(ctx context.Context, postData model.Post, listData model.List) error
	Update(ctx context.Context, data model.List) error
	FindListsByOwner(ctx context.Context, username string) ([]model.List, error)
	FindListByOwnerAndListSlug(ctx context.Context, username, ListSlug string) (*model.List, error)
	FindPostsInAListByListSlug(ctx context.Context, username, ListSlug string) (*model.List, error)
	RemovePostFromList(ctx context.Context, postData model.Post, listData model.List) error
	Delete(ctx context.Context, data model.List) error
}
