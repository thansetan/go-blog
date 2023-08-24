package list

import (
	listhandler "goproject/internal/app/delivery/http/list/handler"
	"goproject/internal/app/delivery/http/middlewares"
	listrepository "goproject/internal/app/repository/list"
	postrepository "goproject/internal/app/repository/post"
	listusecase "goproject/internal/app/usecase/list"
	"log/slog"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Route(r *gin.RouterGroup, db *gorm.DB, logger *slog.Logger) {
	listRepository := listrepository.NewListRepository(db)
	postRepository := postrepository.NewPostRepository(db)
	usecase := listusecase.NewListUsecase(listRepository, postRepository, logger)
	handler := listhandler.NewListHandler(usecase)

	r.POST("/blog/:username/posts/:post_slug/save/:list_slug", middlewares.JWTAuthMiddleware(), handler.AddPostToMyList)
	list := r.Group("/lists/my", middlewares.JWTAuthMiddleware())
	{
		list.POST("", handler.CreateNewList)
		list.GET("", handler.GetMyLists)
		list.GET("/:list_slug", handler.GetPostsInMyListBySlug)
		list.PUT("/:list_slug", handler.UpdateMyListInformationBySlug)
		list.DELETE("/:list_slug", handler.DeleteMyListBySlug)
		list.DELETE("/:list_slug/:post_slug", handler.RemovePostFromMyList)
	}
}
