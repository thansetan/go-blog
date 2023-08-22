package auth

import (
	authhandler "goproject/internal/app/delivery/http/auth/handler"
	blogrepository "goproject/internal/app/repository/blog"
	userrepository "goproject/internal/app/repository/user"
	authusecase "goproject/internal/app/usecase/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Route(r *gin.Engine, db *gorm.DB) {
	userRepository := userrepository.NewUserRepository(db)
	blogRepository := blogrepository.NewBlogRepository(db)
	usecase := authusecase.NewAuthUsecase(userRepository, blogRepository, db)
	handler := authhandler.NewAuthHandler(usecase)

	auth := r.Group("/auth")

	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
	}
}
