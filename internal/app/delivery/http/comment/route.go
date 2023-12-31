package comment

import (
	commenthandler "goproject/internal/app/delivery/http/comment/handler"
	"goproject/internal/app/delivery/http/middlewares"
	commentrepository "goproject/internal/app/repository/comment"
	postrepository "goproject/internal/app/repository/post"
	commentusecase "goproject/internal/app/usecase/comment"
	"log/slog"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Route(r *gin.RouterGroup, db *gorm.DB, logger *slog.Logger) {
	commentRepository := commentrepository.NewCommentRepository(db)
	postRepository := postrepository.NewPostRepository(db)
	usecase := commentusecase.NewCommentUsecase(commentRepository, postRepository, logger)
	handler := commenthandler.NewCommentHandler(usecase)

	r.GET("/my/comments", middlewares.JWTAuthMiddleware(), handler.GetMyComments)
	comment := r.Group("/blog/:username/posts/:post_slug/comments")
	{
		comment.POST("", middlewares.JWTAuthMiddleware(), handler.CreateComment)
		comment.GET("", handler.GetCommentsOnAPost)
		comment.DELETE("/:comment_id", middlewares.JWTAuthMiddleware(), handler.DeleteCommentByID)
		comment.PUT("/:comment_id", middlewares.JWTAuthMiddleware(), handler.EditMyCommentOnAPost)
	}
}
