package middlewares

import (
	"fmt"
	"goproject/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func getToken(raw string) string {
	splitted := strings.Split(raw, " ")
	if len(splitted) > 1 && splitted[0] == "Bearer" {
		return splitted[1]
	}
	return raw
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getToken(c.GetHeader("Authorization"))
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "UNAUTHORIZED",
			})
			c.Abort()
			return
		}

		claims, err := utils.DecodeJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": fmt.Sprintf("Failed to decode JWT: %s", err),
			})
			c.Abort()
			return
		}

		c.Set("username", claims["username"])
		c.Next()
	}
}
