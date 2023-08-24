package postusecase

import (
	"context"
	"errors"
	"fmt"
	"goproject/internal/app/delivery/http/post/dto"
	"goproject/internal/domain/model"
	"goproject/internal/domain/repository"
	"goproject/internal/helpers"
	"log/slog"
	"net/http"

	"gorm.io/gorm"
)

type PostUsecase interface {
	CreateNewPost(ctx context.Context, username string, data dto.PostRequest) (*dto.CreatePostResponse, *helpers.Error)
	GetPostsByBlogOwner(ctx context.Context, username string) ([]dto.PostResponse, *helpers.Error)
	GetPostBySlug(ctx context.Context, username, slug string) (*dto.PostResponse, *helpers.Error)
	UpdatePostBySlug(ctx context.Context, data dto.PostRequest, username, slug string) *helpers.Error
	DeletePostBySlug(ctx context.Context, username, slug string) *helpers.Error
}

type postUsecaseImpl struct {
	postRepo repository.PostRepository
	blogRepo repository.BlogRepository
	logger   *slog.Logger
}

func NewPostUsecase(postRepo repository.PostRepository, blogRepo repository.BlogRepository, logger *slog.Logger) PostUsecase {
	return &postUsecaseImpl{
		postRepo: postRepo,
		blogRepo: blogRepo,
		logger:   logger,
	}
}

func (uc *postUsecaseImpl) CreateNewPost(ctx context.Context, username string, data dto.PostRequest) (*dto.CreatePostResponse, *helpers.Error) {
	resp := new(dto.CreatePostResponse)
	blog, err := uc.blogRepo.FindByOwner(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.ErrorBuilder(http.StatusNotFound, fmt.Sprintf("%s's blog not found", username))
		}
		uc.logger.ErrorContext(ctx, err.Error())
		return nil, helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	postData := model.Post{
		Title:   data.Title,
		Content: data.Content,
		Slug:    helpers.GenerateSlug(data.Title),
		BlogID:  blog.ID,
	}

	err = uc.postRepo.Create(ctx, postData)
	if err != nil {
		uc.logger.ErrorContext(ctx, err.Error())
		return nil, helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	resp.Slug = postData.Slug

	return resp, nil
}

func (uc *postUsecaseImpl) GetPostsByBlogOwner(ctx context.Context, username string) ([]dto.PostResponse, *helpers.Error) {
	postsData := make([]dto.PostResponse, 0)

	blog, err := uc.blogRepo.FindByOwner(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.ErrorBuilder(http.StatusNotFound, fmt.Sprintf("%s's blog not found", username))
		}
		uc.logger.ErrorContext(ctx, err.Error())
		return nil, helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	posts, err := uc.postRepo.FindByBlogID(ctx, blog.ID)
	if err != nil {
		uc.logger.ErrorContext(ctx, err.Error())
		return nil, helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	for _, post := range posts {
		postsData = append(postsData, dto.PostResponse{
			Title:     post.Title,
			Content:   post.Content,
			Slug:      post.Slug,
			Author:    post.Blog.User.Name,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		})
	}

	return postsData, nil
}

func (uc *postUsecaseImpl) GetPostBySlug(ctx context.Context, username, slug string) (*dto.PostResponse, *helpers.Error) {
	post, err := uc.postRepo.FindBySlugAndOwner(ctx, slug, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.ErrorBuilder(http.StatusNotFound, fmt.Sprintf("post %s on %s's blog not found", slug, username))
		}
		uc.logger.ErrorContext(ctx, err.Error())
		return nil, helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}
	data := &dto.PostResponse{
		Title:     post.Title,
		Slug:      post.Slug,
		Content:   post.Content,
		Author:    post.Blog.User.Name,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}
	return data, nil
}

func (uc *postUsecaseImpl) UpdatePostBySlug(ctx context.Context, data dto.PostRequest, username, slug string) *helpers.Error {
	post, err := uc.postRepo.FindBySlugAndOwner(ctx, slug, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helpers.ErrorBuilder(http.StatusNotFound, fmt.Sprintf("post with slug %s not found", slug))
		}
		uc.logger.ErrorContext(ctx, err.Error())
		return helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	post.Title = data.Title
	post.Content = data.Content

	err = uc.postRepo.Update(ctx, *post)
	if err != nil {
		uc.logger.ErrorContext(ctx, err.Error())
		return helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	return nil
}

func (uc *postUsecaseImpl) DeletePostBySlug(ctx context.Context, username, slug string) *helpers.Error {
	post, err := uc.postRepo.FindBySlugAndOwner(ctx, slug, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helpers.ErrorBuilder(http.StatusNotFound, fmt.Sprintf("post with slug %s not found", slug))
		}
		uc.logger.ErrorContext(ctx, err.Error())
		return helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	err = uc.postRepo.Delete(ctx, *post)
	if err != nil {
		uc.logger.ErrorContext(ctx, err.Error())
		return helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	return nil
}
