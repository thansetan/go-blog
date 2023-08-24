package httproute

import (
	"goproject/internal/app/delivery/http/auth"
	"goproject/internal/app/delivery/http/blog"
	"goproject/internal/app/delivery/http/comment"
	"goproject/internal/app/delivery/http/list"
	"goproject/internal/app/delivery/http/post"
	"goproject/internal/app/delivery/http/user"
	"goproject/internal/helpers"

	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func NewRoute(db *gorm.DB, logger *slog.Logger) *gin.Engine {
	r := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("password", helpers.ValidatePassword)
	}

	api := r.Group("/api/v1")

	auth.Route(api, db, logger)
	user.Route(api, db, logger)
	blog.Route(api, db, logger)
	post.Route(api, db, logger)
	comment.Route(api, db, logger)
	list.Route(api, db, logger)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
