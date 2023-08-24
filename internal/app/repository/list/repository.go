package listrepository

import (
	"context"
	"goproject/internal/domain/model"
	"goproject/internal/domain/repository"

	"gorm.io/gorm"
)

type listRepositoryImpl struct {
	db *gorm.DB
}

func NewListRepository(db *gorm.DB) repository.ListRepository {
	return &listRepositoryImpl{
		db: db,
	}
}

func (repo *listRepositoryImpl) Create(ctx context.Context, data model.List) error {
	err := repo.db.WithContext(ctx).Create(&data).Error
	return err
}

func (repo *listRepositoryImpl) Update(ctx context.Context, data model.List) error {
	newData := map[string]any{
		"name":        data.Name,
		"description": data.Description,
	}

	err := repo.db.WithContext(ctx).Model(data).Updates(newData).Error
	return err
}

func (repo *listRepositoryImpl) FindListsByOwner(ctx context.Context, username string) ([]model.List, error) {
	var lists []model.List
	err := repo.db.WithContext(ctx).Preload("Posts").Find(&lists, "owner = ?", username).Error
	if err != nil {
		return nil, err
	}
	return lists, nil
}

func (repo *listRepositoryImpl) FindListByOwnerAndListSlug(ctx context.Context, username string, listSlug string) (*model.List, error) {
	list := new(model.List)
	err := repo.db.WithContext(ctx).First(&list, "owner=? AND slug=?", username, listSlug).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (repo *listRepositoryImpl) FindPostsInAListByListSlug(ctx context.Context, username string, listSlug string) (*model.List, error) {
	list := new(model.List)
	err := repo.db.WithContext(ctx).Preload("Posts.Blog.User").Preload("Posts").First(&list, "owner=? AND slug=?", username, listSlug).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (repo *listRepositoryImpl) AddPostToList(ctx context.Context, postData model.Post, listData model.List) error {
	// err := repo.db.WithContext(ctx).Model(&listData).Association("Posts").Append(&postData) // this won't throw an error on duplicate posts in the same list
	err := repo.db.WithContext(ctx).Create(&model.ListPost{
		ListID: listData.ID,
		PostID: postData.ID,
	}).Error
	return err
}

func (repo *listRepositoryImpl) RemovePostFromList(ctx context.Context, postData model.Post, listData model.List) error {
	err := repo.db.WithContext(ctx).Model(&listData).Association("Posts").Delete(&postData)
	return err
}

func (repo *listRepositoryImpl) Delete(ctx context.Context, data model.List) error {
	tx := repo.db.Begin()
	err := tx.WithContext(ctx).Model(&data).Association("ListPosts").Unscoped().Clear()
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.WithContext(ctx).Delete(&data).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
