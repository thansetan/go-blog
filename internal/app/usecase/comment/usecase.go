package commentusecase

import (
	"context"
	"errors"
	"fmt"
	"goproject/internal/app/delivery/http/comment/dto"
	"goproject/internal/domain/model"
	"goproject/internal/domain/repository"
	"goproject/internal/helpers"
	"log/slog"
	"net/http"
	"slices"
	"strconv"

	"gorm.io/gorm"
)

type CommentUsecase interface {
	CreateComment(ctx context.Context, data dto.CommentRequest, username, blogOwner, postSlug string) (uint, *helpers.Error)
	GetCommentsByUsername(ctx context.Context, username string) ([]dto.CommentResponse, *helpers.Error)
	GetCommentsByBlogOwnerAndPostSlug(ctx context.Context, blogOwner, postSlug string) ([]dto.CommentResponse, *helpers.Error)
	DeleteCommentOnAPosst(ctx context.Context, username, blogOwner, postSlug string, commentID string) *helpers.Error
	UpdateCommentOnAPost(ctx context.Context, username, blogOwner, postSlug string, commentID string, data dto.CommentRequest) *helpers.Error
}

type CommentUsecaseImpl struct {
	commentRepo repository.CommentRepository
	postRepo    repository.PostRepository
	logger      *slog.Logger
}

func NewCommentUsecase(commentRepo repository.CommentRepository, postRepo repository.PostRepository, logger *slog.Logger) CommentUsecase {
	return &CommentUsecaseImpl{
		postRepo:    postRepo,
		commentRepo: commentRepo,
		logger:      logger,
	}
}

func (uc *CommentUsecaseImpl) CreateComment(ctx context.Context, data dto.CommentRequest, username, blogOwner, postSlug string) (uint, *helpers.Error) {
	post, err := uc.postRepo.FindBySlugAndOwner(ctx, postSlug, blogOwner)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, helpers.ErrorBuilder(http.StatusNotFound, fmt.Sprintf("post %s on %s's blog not found", postSlug, blogOwner))
		}
		uc.logger.ErrorContext(ctx, err.Error())
		return 0, helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	commentData := model.Comment{
		Commenter: username,
		PostID:    post.ID,
		Content:   data.Comment,
	}

	id, err := uc.commentRepo.Create(ctx, commentData)
	if err != nil {
		uc.logger.ErrorContext(ctx, err.Error())
		return 0, helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}
	return id, nil
}

func (uc *CommentUsecaseImpl) GetCommentsByUsername(ctx context.Context, username string) ([]dto.CommentResponse, *helpers.Error) {
	commentsData := make([]dto.CommentResponse, 0)

	comments, err := uc.commentRepo.FindCommentByUsername(ctx, username)
	if err != nil {
		uc.logger.ErrorContext(ctx, err.Error())
		return nil, helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
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

func (uc *CommentUsecaseImpl) GetCommentsByBlogOwnerAndPostSlug(ctx context.Context, blogOwner, postSlug string) ([]dto.CommentResponse, *helpers.Error) {
	commentsData := make([]dto.CommentResponse, 0)

	post, err := uc.postRepo.FindBySlugAndOwner(ctx, postSlug, blogOwner)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.ErrorBuilder(http.StatusNotFound, fmt.Sprintf("post %s on %s's blog not found", postSlug, blogOwner))
		}
		uc.logger.ErrorContext(ctx, err.Error())
		return nil, helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	comments, err := uc.commentRepo.FindCommentByPostID(ctx, post.ID)
	if err != nil {
		uc.logger.ErrorContext(ctx, err.Error())
		return nil, helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
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

func (uc *CommentUsecaseImpl) DeleteCommentOnAPosst(ctx context.Context, username, blogOwner, postSlug string, commentID string) *helpers.Error {
	// blog owner can delete any users comment on their post
	// user (non blog owner) can ony delete their own comments on someone else's blog post

	post, err := uc.postRepo.FindBySlugAndOwner(ctx, postSlug, blogOwner)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helpers.ErrorBuilder(http.StatusNotFound, fmt.Sprintf("post %s on %s's blog not found", postSlug, blogOwner))
		}
		uc.logger.ErrorContext(ctx, err.Error())
		return helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	comments, err := uc.commentRepo.FindCommentByPostID(ctx, post.ID)
	if err != nil {
		uc.logger.ErrorContext(ctx, err.Error())
		return helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	id, err := strconv.ParseUint(commentID, 10, 64)
	if err != nil {
		return helpers.ErrorBuilder(http.StatusBadRequest, "comment id must be a positive integer")
	}

	commentIdx := slices.IndexFunc(comments, func(c model.Comment) bool {
		return c.ID == uint(id)
	})

	if commentIdx == -1 {
		return helpers.ErrorBuilder(http.StatusNotFound, fmt.Sprintf("comment with id %s not found on this post", commentID))
	}

	if comments[commentIdx].Commenter != username && username != blogOwner {
		return helpers.ErrorBuilder(http.StatusUnauthorized, "you're not allowed to delete this comment")
	}

	err = uc.commentRepo.Delete(ctx, comments[commentIdx])
	if err != nil {
		uc.logger.ErrorContext(ctx, err.Error())
		return helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	return nil
}

func (uc *CommentUsecaseImpl) UpdateCommentOnAPost(ctx context.Context, username, blogOwner, postSlug string, commentID string, data dto.CommentRequest) *helpers.Error {
	// the only person able to edit a comment is the commenter

	post, err := uc.postRepo.FindBySlugAndOwner(ctx, postSlug, blogOwner)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helpers.ErrorBuilder(http.StatusNotFound, fmt.Sprintf("post %s on %s's blog not found", postSlug, blogOwner))
		}
		uc.logger.ErrorContext(ctx, err.Error())
		return helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	comments, err := uc.commentRepo.FindCommentByPostID(ctx, post.ID)
	if err != nil {
		uc.logger.ErrorContext(ctx, err.Error())
		return helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	id, err := strconv.ParseUint(commentID, 10, 64)
	if err != nil {
		return helpers.ErrorBuilder(http.StatusBadRequest, "comment id must be a positive integer")
	}

	commentIdx := slices.IndexFunc(comments, func(c model.Comment) bool {
		return c.ID == uint(id)
	})

	if commentIdx == -1 {
		return helpers.ErrorBuilder(http.StatusNotFound, fmt.Sprintf("comment with id %s not found on this post", commentID))
	}

	if comments[commentIdx].Commenter != username {
		return helpers.ErrorBuilder(http.StatusUnauthorized, "you're not allowed to modify this comment")
	}

	comments[commentIdx].Content = data.Comment

	err = uc.commentRepo.Update(ctx, comments[commentIdx])
	if err != nil {
		uc.logger.ErrorContext(ctx, err.Error())
		return helpers.ErrorBuilder(http.StatusInternalServerError, "it's our fault, not yours")
	}

	return nil
}
