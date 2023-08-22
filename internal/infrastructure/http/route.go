package httproute

import (
	"goproject/internal/app/delivery/http/auth"
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

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// r.Use(middlewares.JWTAuthMiddleware()).GET("/protected", func(c *gin.Context) {
	// 	c.String(http.StatusOK, "KEREN")
	// })

	return r
}
