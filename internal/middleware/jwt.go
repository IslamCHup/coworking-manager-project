package middleware

import (
	"fmt"
	"strings"

	authjwt "github.com/IslamCHup/coworking-manager-project/internal/auth/jwt"
	"github.com/gin-gonic/gin"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.Next()
			return
		}

		parts := strings.Split(header, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		claims, err := authjwt.ParseAccessToken(parts[1])
		if err == nil && claims != nil {
			c.Set("user_id", claims.UserID)
		} else {
			fmt.Println("JWT MIDDLEWARE: invalid token:", err)
		}

		c.Next()
	}
}
