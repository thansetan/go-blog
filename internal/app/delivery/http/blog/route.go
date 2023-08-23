package blog

import (
	bloghandler "goproject/internal/app/delivery/http/blog/handler"
	"goproject/internal/app/delivery/http/middlewares"
	blogrepository "goproject/internal/app/repository/blog"
	blogusecase "goproject/internal/app/usecase/blog"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Route(r *gin.RouterGroup, db *gorm.DB) {
	repository := blogrepository.NewBlogRepository(db)
	usecase := blogusecase.NewBlogUsecase(repository)
	handler := bloghandler.NewBlogHandler(usecase)

	blog := r.Group("/blog")
	{
		myBlog := blog.Group("/my").Use(middlewares.JWTAuthMiddleware())
		{
			myBlog.PUT("", handler.UpdateMyBlogName)
			myBlog.GET("", handler.GetMyBlog)
		}
		blog.GET("/:username", handler.GetBlogByOwner)
	}
}
