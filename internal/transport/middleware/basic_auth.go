package middleware

import (
	"encoding/base64"
	"log/slog"
	"net/http"
	"strings"

	"github.com/IslamCHup/coworking-manager-project/internal/service"
	"github.com/gin-gonic/gin"
)

// AdminBasicAuthMiddleware создаёт middleware для HTTP Basic Authentication
// Проверяет логин и пароль через AdminService
func AdminBasicAuthMiddleware(adminService service.AdminService, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Header("WWW-Authenticate", `Basic realm="admin"`)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization required"})
			c.Abort()
			return
		}

		// Проверяем, что заголовок начинается с "Basic "
		if !strings.HasPrefix(authHeader, "Basic ") {
			c.Header("WWW-Authenticate", `Basic realm="admin"`)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		// Извлекаем base64-кодированную часть
		encoded := strings.TrimPrefix(authHeader, "Basic ")
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			c.Header("WWW-Authenticate", `Basic realm="admin"`)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization credentials"})
			c.Abort()
			return
		}

		// Декодированная строка должна быть в формате "username:password"
		credentials := strings.SplitN(string(decoded), ":", 2)
		if len(credentials) != 2 {
			c.Header("WWW-Authenticate", `Basic realm="admin"`)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization credentials format"})
			c.Abort()
			return
		}

		login := credentials[0]
		password := credentials[1]

		// Проверяем через сервис (из БД)
		admin, err := adminService.VerifyAdmin(login, password)
		if err != nil {
			logger.Warn("Admin authentication failed", "login", login, "error", err)
			c.Header("WWW-Authenticate", `Basic realm="admin"`)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid login or password"})
			c.Abort()
			return
		}

		// Сохраняем админа в контексте для использования в хендлерах
		c.Set("admin_id", admin.ID)
		c.Set("admin_login", admin.Login)

		// Аутентификация успешна, продолжаем выполнение
		c.Next()
	}
}
