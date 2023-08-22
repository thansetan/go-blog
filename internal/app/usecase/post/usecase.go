package postusecase

import (
	"context"
	"goproject/internal/app/delivery/http/post/dto"
	"goproject/internal/domain/model"
	"goproject/internal/domain/repository"
	"goproject/internal/utils"
)

type PostUsecase interface {
	CreateNewPost(ctx context.Context, username string, data dto.PostRequest) error
	GetPostsByBlogOwner(ctx context.Context, username string) ([]dto.PostResponse, error)
	GetPostBySlug(ctx context.Context, username, slug string) (*dto.PostResponse, error)
	UpdatePostBySlug(ctx context.Context, data dto.PostRequest, username, slug string) error
	DeletePostBySlug(ctx context.Context, username, slug string) error
}

type PostUsecaseImpl struct {
	postRepo repository.PostRepository
	blogRepo repository.BlogRepository
}

func NewPostUsecase(postRepo repository.PostRepository, blogRepo repository.BlogRepository) PostUsecase {
	return &PostUsecaseImpl{
		postRepo: postRepo,
		blogRepo: blogRepo,
	}
}

func (uc *PostUsecaseImpl) CreateNewPost(ctx context.Context, username string, data dto.PostRequest) error {
	blog, err := uc.blogRepo.FindByOwner(ctx, username)
	if err != nil {
		return err
	}

	titleSlug := utils.GenerateSlug(data.Title)
	postData := model.Post{
		Title:     data.Title,
		Content:   data.Content,
		TitleSlug: titleSlug,
		Author:    blog.User.Name,
		BlogID:    blog.ID,
	}

	err = uc.postRepo.Create(ctx, postData)
	if err != nil {
		return err
	}

	return nil
}

func (uc *PostUsecaseImpl) GetPostsByBlogOwner(ctx context.Context, username string) ([]dto.PostResponse, error) {
	var postsData []dto.PostResponse

	blog, err := uc.blogRepo.FindByOwner(ctx, username)
	if err != nil {
		return nil, err
	}

	posts, err := uc.postRepo.FindByBlogID(ctx, blog.ID)
	if err != nil {
		return nil, err
	}

	for _, post := range posts {
		postsData = append(postsData, dto.PostResponse{
			Title:     post.Title,
			Content:   post.Content,
			Slug:      post.TitleSlug,
			Author:    post.Author,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		})
	}

	return postsData, nil
}

func (uc *PostUsecaseImpl) GetPostBySlug(ctx context.Context, username, slug string) (*dto.PostResponse, error) {
	post, err := uc.postRepo.FindBySlugAndOwner(ctx, slug, username)
	if err != nil {
		return nil, err
	}
	data := &dto.PostResponse{
		Title:     post.Title,
		Slug:      post.TitleSlug,
		Content:   post.Content,
		Author:    post.Author,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}
	return data, nil
}

func (uc *PostUsecaseImpl) UpdatePostBySlug(ctx context.Context, data dto.PostRequest, username, slug string) error {
	post, err := uc.postRepo.FindBySlugAndOwner(ctx, slug, username)
	if err != nil {
		return err
	}

	post.Title = data.Title
	post.Content = data.Content

	err = uc.postRepo.Update(ctx, *post)
	if err != nil {
		return err
	}

	return nil
}

func (uc *PostUsecaseImpl) DeletePostBySlug(ctx context.Context, username, slug string) error {
	post, err := uc.postRepo.FindBySlugAndOwner(ctx, slug, username)
	if err != nil {
		return err
	}

	err = uc.postRepo.Delete(ctx, *post)
	if err != nil {
		return err
	}

	return nil
}
