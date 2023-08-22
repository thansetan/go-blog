package comment

import (
	commenthandler "goproject/internal/app/delivery/http/comment/handler"
	"goproject/internal/app/delivery/http/middlewares"
	commentrepository "goproject/internal/app/repository/comment"
	postrepository "goproject/internal/app/repository/post"
	commentusecase "goproject/internal/app/usecase/comment"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Route(r *gin.Engine, db *gorm.DB) {
	commentRepository := commentrepository.NewCommentRepository(db)
	postRepository := postrepository.NewPostRepository(db)
	usecase := commentusecase.NewCommentUsecase(commentRepository, postRepository)
	handler := commenthandler.NewCommentHandler(usecase)

	r.GET("/my/comments", middlewares.JWTAuthMiddleware(), handler.GetMyComments)
	comment := r.Group("/blog/:username/posts/:slug/comments")
	{
		comment.POST("", middlewares.JWTAuthMiddleware(), handler.CreateComment)
		comment.GET("", handler.GetCommentsOnAPost)
		comment.DELETE("/:comment_id", middlewares.JWTAuthMiddleware(), handler.DeleteCommentByID)
		comment.PUT("/:comment_id", middlewares.JWTAuthMiddleware(), handler.EditMyCommentOnAPost)
	}
}
