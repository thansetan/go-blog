package listusecase

import (
	"context"
	"errors"
	"fmt"
	"goproject/internal/app/delivery/http/list/dto"
	postDto "goproject/internal/app/delivery/http/post/dto"
	"goproject/internal/domain/model"
	"goproject/internal/domain/repository"
	"goproject/internal/helpers"
	"log/slog"
	"net/http"
	"slices"

	"gorm.io/gorm"
)

type ListUsecase interface {
	CreateNewList(ctx context.Context, data dto.ListRequest, username string) (*dto.CreateListResponse, *helpers.Error)
	AddPostToMyList(ctx context.Context, listSlug, username, blogOwner, postSlug string) *helpers.Error
	GetPostsInAListBySlug(ctx context.Context, listSlug, username string) (*dto.ListResponse, *helpers.Error)
	GetMyLists(ctx context.Context, username string) ([]dto.ListResponse, *helpers.Error)
	UpdateListInformation(ctx context.Context, data dto.ListRequest, username, listSlug string) *helpers.Error
	RemovePostFromList(ctx context.Context, username, slug, listSlug string) *helpers.Error
	DeleteListBySlug(ctx context.Context, username, listSlug string) *helpers.Error
}

type listUsecaseImpl struct {
	listRepo repository.ListRepository
	postRepo repository.PostRepository
	logger   *slog.Logger
}

func NewListUsecase(listRepo repository.ListRepository, postRepo repository.PostRepository, logger *slog.Logger) ListUsecase {
	return &listUsecaseImpl{
		listRepo: listRepo,
		postRepo: postRepo,
		logger:   logger,
	}
}

func (uc *listUsecaseImpl) CreateNewList(ctx context.Context, data dto.ListRequest, username string) (*dto.CreateListResponse, *helpers.Error) {
	resp := new(dto.CreateListResponse)

	listData := model.List{
		Name:        data.Name,
		Slug:        helpers.GenerateSlug(data.Name),
		Description: data.Description,
		Owner:       username,
	}

	err := uc.listRepo.Create(ctx, listData)
	if err != nil {
		uc.logger.ErrorContext(ctx, err.Error())
		return nil, helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	resp.Slug = listData.Slug

	return resp, nil
}

func (uc *listUsecaseImpl) AddPostToMyList(ctx context.Context, listSlug, username, blogOwner, postSlug string) *helpers.Error {
	post, err := uc.postRepo.FindBySlugAndOwner(ctx, postSlug, blogOwner)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helpers.ErrorBuilder(http.StatusNotFound, fmt.Sprintf("post %s on %s's blog not found", postSlug, blogOwner))
		}
		uc.logger.ErrorContext(ctx, err.Error())
		return helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	list, err := uc.listRepo.FindListByOwnerAndListSlug(ctx, username, listSlug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helpers.ErrorBuilder(http.StatusNotFound, fmt.Sprintf("list %s not found", listSlug))
		}
		uc.logger.ErrorContext(ctx, err.Error())
		return helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	err = uc.listRepo.AddPostToList(ctx, *post, *list)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return helpers.ErrorBuilder(http.StatusConflict, "this post is already in this list!")
		}
		uc.logger.ErrorContext(ctx, err.Error())
		return helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	return nil
}

func (uc *listUsecaseImpl) GetPostsInAListBySlug(ctx context.Context, listSlug, username string) (*dto.ListResponse, *helpers.Error) {
	listData := new(dto.ListResponse)
	listData.Posts = &[]postDto.PostResponse{}

	list, err := uc.listRepo.FindPostsInAListByListSlug(ctx, username, listSlug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.ErrorBuilder(http.StatusNotFound, fmt.Sprintf("list %s not found", listSlug))
		}
		uc.logger.ErrorContext(ctx, err.Error())
		return nil, helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	listData.Slug = list.Slug
	listData.Name = list.Name
	listData.Description = list.Description

	for _, post := range list.Posts {
		*listData.Posts = append(*listData.Posts, postDto.PostResponse{
			Title:     post.Title,
			Content:   post.Content,
			Slug:      post.Slug,
			Author:    post.Blog.User.Name,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		})
	}

	return listData, nil
}

func (uc *listUsecaseImpl) GetMyLists(ctx context.Context, username string) ([]dto.ListResponse, *helpers.Error) {
	listsData := make([]dto.ListResponse, 0)

	lists, err := uc.listRepo.FindListsByOwner(ctx, username)
	if err != nil {
		uc.logger.ErrorContext(ctx, err.Error())
		return nil, helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
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

func (uc *listUsecaseImpl) UpdateListInformation(ctx context.Context, data dto.ListRequest, username, listSlug string) *helpers.Error {
	list, err := uc.listRepo.FindListByOwnerAndListSlug(ctx, username, listSlug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helpers.ErrorBuilder(http.StatusNotFound, fmt.Sprintf("list %s not found", listSlug))
		}
		uc.logger.ErrorContext(ctx, err.Error())
		return helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	list.Name = data.Name
	list.Description = data.Description

	err = uc.listRepo.Update(ctx, *list)
	if err != nil {
		uc.logger.ErrorContext(ctx, err.Error())
		return helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}
	return nil
}

func (uc *listUsecaseImpl) RemovePostFromList(ctx context.Context, username, postSlug, listSlug string) *helpers.Error {
	list, err := uc.listRepo.FindPostsInAListByListSlug(ctx, username, listSlug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helpers.ErrorBuilder(http.StatusNotFound, fmt.Sprintf("list %s not found", listSlug))
		}
		uc.logger.ErrorContext(ctx, err.Error())
		return helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	postIdx := slices.IndexFunc(list.Posts, func(post model.Post) bool {
		return post.Slug == postSlug
	})
	if postIdx == -1 {
		return helpers.ErrorBuilder(http.StatusNotFound, fmt.Sprintf("post %s not found on this list", postSlug))
	}

	post := list.Posts[postIdx]

	err = uc.listRepo.RemovePostFromList(ctx, post, *list)
	if err != nil {
		uc.logger.ErrorContext(ctx, err.Error())
		return helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	return nil
}

func (uc *listUsecaseImpl) DeleteListBySlug(ctx context.Context, username, listSlug string) *helpers.Error {
	list, err := uc.listRepo.FindListByOwnerAndListSlug(ctx, username, listSlug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helpers.ErrorBuilder(http.StatusNotFound, fmt.Sprintf("list %s not found", listSlug))
		}
		uc.logger.ErrorContext(ctx, err.Error())
		return helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	err = uc.listRepo.Delete(ctx, *list)
	if err != nil {
		uc.logger.ErrorContext(ctx, err.Error())
		return helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	return nil
}
