package user

import (
	"goproject/internal/app/delivery/http/middlewares"
	userhandler "goproject/internal/app/delivery/http/user/handler"
	userrepository "goproject/internal/app/repository/user"
	userusecase "goproject/internal/app/usecase/user"

	"log/slog"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Route(r *gin.RouterGroup, db *gorm.DB, logger *slog.Logger) {
	repository := userrepository.NewUserRepository(db)
	usecase := userusecase.NewUserUsecase(repository, logger)
	handler := userhandler.NewUserHandler(usecase)

	user := r.Group("/users")
	user.Use(middlewares.JWTAuthMiddleware())
	{
		user.GET("/me", handler.GetMyInformation)
		user.PUT("/me", handler.UpdateMyInformation)
		user.PUT("/me/update-password", handler.UpdateMyPassword)

	}
}
