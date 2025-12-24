package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	authjwt "github.com/IslamCHup/coworking-manager-project/internal/auth/jwt"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "отсутствует заголовок авторизации",
			})
			return
		}

		parts := strings.Split(header, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "неверный формат авторизации",
			})
			return
		}

		claims, err := authjwt.ParseAccessToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "неверный или истекший токен",
			})
			return
		}

		c.Set("userID", claims.UserID)

		c.Next()
	}
}
