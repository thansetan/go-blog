package httproute

import (
	"goproject/internal/app/delivery/http/auth"
	"goproject/internal/app/delivery/http/blog"
	"goproject/internal/app/delivery/http/comment"
	"goproject/internal/app/delivery/http/post"
	"goproject/internal/app/delivery/http/user"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func NewRoute(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	auth.Route(r, db)
	user.Route(r, db)
	blog.Route(r, db)
	post.Route(r, db)
	comment.Route(r, db)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
