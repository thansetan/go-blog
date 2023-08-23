package httproute

import (
	"goproject/internal/app/delivery/http/auth"
	"goproject/internal/app/delivery/http/blog"
	"goproject/internal/app/delivery/http/comment"
	"goproject/internal/app/delivery/http/list"
	"goproject/internal/app/delivery/http/post"
	"goproject/internal/app/delivery/http/user"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func NewRoute(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	api := r.Group("/api/v1")

	auth.Route(api, db)
	user.Route(api, db)
	blog.Route(api, db)
	post.Route(api, db)
	comment.Route(api, db)
	list.Route(api, db)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
