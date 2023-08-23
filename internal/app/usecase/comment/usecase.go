package commentusecase

import (
	"context"
	"errors"
	"fmt"
	"goproject/internal/app/delivery/http/comment/dto"
	"goproject/internal/domain/model"
	"goproject/internal/domain/repository"
	"slices"
	"strconv"
)

type CommentUsecase interface {
	CreateComment(ctx context.Context, data dto.CommentRequest, username, blogOwner, postSlug string) error
	GetCommentByUsername(ctx context.Context, username string) ([]dto.CommentResponse, error)
	GetCommentByBlogOwnerAndPostSlug(ctx context.Context, blogOwner, postSlug string) ([]dto.CommentResponse, error)
	DeleteCommentOnAPosst(ctx context.Context, username, blogOwner, postSlug string, commentID string) error
	UpdateCommentOnAPost(ctx context.Context, username, blogOwner, postSlug string, commentID string, data dto.CommentRequest) error
}

type CommentUsecaseImpl struct {
	commentRepo repository.CommentRepository
	postRepo    repository.PostRepository
}

func NewCommentUsecase(commentRepo repository.CommentRepository, postRepo repository.PostRepository) CommentUsecase {
	return &CommentUsecaseImpl{
		postRepo:    postRepo,
		commentRepo: commentRepo,
	}
}

func (uc *CommentUsecaseImpl) CreateComment(ctx context.Context, data dto.CommentRequest, username, blogOwner, postSlug string) error {
	post, err := uc.postRepo.FindBySlugAndOwner(ctx, postSlug, blogOwner)
	if err != nil {
		return err
	}

	commentData := model.Comment{
		Commenter: username,
		PostID:    post.ID,
		Content:   data.Comment,
	}

	err = uc.commentRepo.Create(ctx, commentData)
	if err != nil {
		return err
	}
	return nil
}

func (uc *CommentUsecaseImpl) GetCommentByUsername(ctx context.Context, username string) ([]dto.CommentResponse, error) {
	var commentsData []dto.CommentResponse

	comments, err := uc.commentRepo.FindCommentByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	for _, comment := range comments {
		commentsData = append(commentsData, dto.CommentResponse{
			ID:        comment.ID,
			PostURL:   fmt.Sprintf("blog/%s/posts/%s", comment.Post.Blog.Owner, comment.Post.Slug),
			Comment:   comment.Content,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
		})
	}

	return commentsData, nil
}

func (uc *CommentUsecaseImpl) GetCommentByBlogOwnerAndPostSlug(ctx context.Context, blogOwner, postSlug string) ([]dto.CommentResponse, error) {
	var commentsData []dto.CommentResponse

	post, err := uc.postRepo.FindBySlugAndOwner(ctx, postSlug, blogOwner)
	if err != nil {
		return nil, err
	}

	comments, err := uc.commentRepo.FindCommentByPostID(ctx, post.ID)
	if err != nil {
		return nil, err
	}

	for _, comment := range comments {
		commentsData = append(commentsData, dto.CommentResponse{
			ID:        comment.ID,
			Commenter: comment.Commenter,
			Comment:   comment.Content,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
		})
	}

	return commentsData, nil
}

func (uc *CommentUsecaseImpl) DeleteCommentOnAPosst(ctx context.Context, username, blogOwner, postSlug string, commentID string) error {
	// blog owner can delete any users comment on their post
	// user (non blog owner) can ony delete their own comments on someone else's blog post

	post, err := uc.postRepo.FindBySlugAndOwner(ctx, postSlug, blogOwner)
	if err != nil {
		return err
	}

	comments, err := uc.commentRepo.FindCommentByPostID(ctx, post.ID)
	if err != nil {
		return err
	}

	id, err := strconv.ParseUint(commentID, 10, 64)
	if err != nil {
		return err
	}

	commentIdx := slices.IndexFunc(comments, func(c model.Comment) bool {
		return c.ID == uint(id)
	})

	if commentIdx == -1 {
		return errors.New("comment not in this posts")
	}

	if comments[commentIdx].Commenter != username && username != blogOwner {
		fmt.Println(comments[commentIdx].Commenter, username, blogOwner)
		return errors.New("you're not allowed to delete this comment")
	}

	err = uc.commentRepo.Delete(ctx, comments[commentIdx])
	if err != nil {
		return err
	}

	return nil
}

func (uc *CommentUsecaseImpl) UpdateCommentOnAPost(ctx context.Context, username, blogOwner, postSlug string, commentID string, data dto.CommentRequest) error {
	// the only person able to edit a comment is the commenter

	post, err := uc.postRepo.FindBySlugAndOwner(ctx, postSlug, blogOwner)
	if err != nil {
		return err
	}

	comments, err := uc.commentRepo.FindCommentByPostID(ctx, post.ID)
	if err != nil {
		return err
	}

	id, err := strconv.ParseUint(commentID, 10, 64)
	if err != nil {
		return err
	}

	commentIdx := slices.IndexFunc(comments, func(c model.Comment) bool {
		return c.ID == uint(id)
	})

	if commentIdx == -1 {
		return errors.New("comment not in this posts")
	}

	if comments[commentIdx].Commenter != username {
		return errors.New("you're not allowed to modify this comment")
	}

	comments[commentIdx].Content = data.Comment

	err = uc.commentRepo.Update(ctx, comments[commentIdx])
	if err != nil {
		return err
	}

	return nil
}
