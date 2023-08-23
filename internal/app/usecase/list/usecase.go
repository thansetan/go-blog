package listusecase

import (
	"context"
	"errors"
	"fmt"
	"goproject/internal/app/delivery/http/list/dto"
	postDto "goproject/internal/app/delivery/http/post/dto"
	"goproject/internal/domain/model"
	"goproject/internal/domain/repository"
	"goproject/internal/utils"
	"slices"
)

type ListUsecase interface {
	CreateNewList(ctx context.Context, data dto.ListRequest, username string) error
	AddPostToMyList(ctx context.Context, listSlug, username, blogOwner, postSlug string) error
	GetPostsInAListBySlug(ctx context.Context, listSlug, username string) (*dto.ListResponse, error)
	GetMyLists(ctx context.Context, username string) ([]dto.ListResponse, error)
	UpdateListInformation(ctx context.Context, data dto.ListRequest, username, listSlug string) error
	RemovePostFromList(ctx context.Context, username, slug, listSlug string) error
	DeleteListBySlug(ctx context.Context, username, listSlug string) error
}

type ListUsecaseImpl struct {
	listRepo repository.ListRepository
	postRepo repository.PostRepository
}

func NewListUsecase(listRepo repository.ListRepository, postRepo repository.PostRepository) ListUsecase {
	return &ListUsecaseImpl{
		listRepo: listRepo,
		postRepo: postRepo,
	}
}

func (uc *ListUsecaseImpl) CreateNewList(ctx context.Context, data dto.ListRequest, username string) error {
	listData := model.List{
		Name:        data.Name,
		Slug:        utils.GenerateSlug(data.Name),
		Description: data.Description,
		Owner:       username,
	}

	err := uc.listRepo.Create(ctx, listData)
	if err != nil {
		return err
	}
	return nil
}

func (uc *ListUsecaseImpl) AddPostToMyList(ctx context.Context, listSlug, username, blogOwner, postSlug string) error {
	post, err := uc.postRepo.FindBySlugAndOwner(ctx, postSlug, blogOwner)
	if err != nil {
		return err
	}

	list, err := uc.listRepo.FindListByOwnerAndListSlug(ctx, username, listSlug)
	if err != nil {
		return err
	}

	err = uc.listRepo.AddPostToList(ctx, *post, *list)
	if err != nil {
		return err
	}

	return nil
}

func (uc *ListUsecaseImpl) GetPostsInAListBySlug(ctx context.Context, listSlug, username string) (*dto.ListResponse, error) {
	listData := new(dto.ListResponse)
	listData.Posts = &[]postDto.PostResponse{}

	list, err := uc.listRepo.FindPostsInAListByListSlug(ctx, username, listSlug)
	if err != nil {
		return nil, err
	}

	listData.Slug = list.Slug
	listData.Name = list.Name
	listData.Description = list.Description

	for _, post := range list.Posts {
		*listData.Posts = append(*listData.Posts, postDto.PostResponse{
			Title:     post.Title,
			Content:   post.Content,
			Slug:      post.Slug,
			Author:    post.Author,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		})
	}
	fmt.Println(listData.Posts == nil)
	fmt.Println(listData)
	return listData, nil
}

func (uc *ListUsecaseImpl) GetMyLists(ctx context.Context, username string) ([]dto.ListResponse, error) {
	listsData := make([]dto.ListResponse, 0)

	lists, err := uc.listRepo.FindListsByOwner(ctx, username)
	if err != nil {
		return nil, err
	}

	for _, list := range lists {
		numOfPosts := len(list.Posts)
		listsData = append(listsData, dto.ListResponse{
			Slug:        list.Slug,
			Name:        list.Name,
			Description: list.Description,
			NumOfPosts:  &numOfPosts,
		})
	}

	return listsData, nil
}

func (uc *ListUsecaseImpl) UpdateListInformation(ctx context.Context, data dto.ListRequest, username, listSlug string) error {
	list, err := uc.listRepo.FindListByOwnerAndListSlug(ctx, username, listSlug)
	if err != nil {
		return err
	}

	list.Name = data.Name
	list.Description = data.Description

	err = uc.listRepo.Update(ctx, *list)
	if err != nil {
		return err
	}
	return nil
}

func (uc *ListUsecaseImpl) RemovePostFromList(ctx context.Context, username, slug, listSlug string) error {
	list, err := uc.listRepo.FindPostsInAListByListSlug(ctx, username, listSlug)
	if err != nil {
		return err
	}

	postIdx := slices.IndexFunc(list.Posts, func(post model.Post) bool {
		return post.Slug == slug
	})
	if postIdx == -1 {
		return errors.New("post not found in this list")
	}

	post := list.Posts[postIdx]

	err = uc.listRepo.RemovePostFromList(ctx, post, *list)
	if err != nil {
		return err
	}

	return nil
}

func (uc *ListUsecaseImpl) DeleteListBySlug(ctx context.Context, username, listSlug string) error {
	list, err := uc.listRepo.FindListByOwnerAndListSlug(ctx, username, listSlug)
	if err != nil {
		fmt.Println("ERROR DARI SINI 2")
		return err
	}

	err = uc.listRepo.Delete(ctx, *list)
	if err != nil {
		fmt.Println("ERROR DARI SINI 3")
		return err
	}

	return nil
}
