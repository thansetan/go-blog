package post

import (
	"goproject/internal/app/delivery/http/middlewares"
	posthandler "goproject/internal/app/delivery/http/post/handler"
	blogrepository "goproject/internal/app/repository/blog"
	postrepository "goproject/internal/app/repository/post"
	postusecase "goproject/internal/app/usecase/post"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Route(r *gin.Engine, db *gorm.DB) {
	postRepository := postrepository.NewPostRepository(db)
	blogRepository := blogrepository.NewBlogRepository(db)
	usecase := postusecase.NewPostUsecase(postRepository, blogRepository)
	handler := posthandler.NewPostHandler(usecase)

	r.GET("/blog/:username/posts", handler.GetPostsByBlogOwner)
	r.GET("/blog/:username/posts/:post_slug", handler.GetPostBySlug)

	post := r.Group("/blog/my/posts", middlewares.JWTAuthMiddleware())
	{
		post.POST("", handler.CreateNewPost)
		post.GET("", handler.GetAllMyBlogPosts)
		post.PUT("/:post_slug", handler.UpdateMyPostBySlug)
		post.DELETE("/:post_slug", handler.DeleteMyPostBySlug)
	}
}
