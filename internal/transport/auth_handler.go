package transport

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/IslamCHup/coworking-manager-project/internal/auth/jwt"
	"github.com/IslamCHup/coworking-manager-project/internal/models"
	"github.com/IslamCHup/coworking-manager-project/internal/service"
)

type AuthHandler struct {
	authService    service.AuthService
	refreshService service.RefreshService
	logger         *slog.Logger
}

func NewAuthHandler(
	authService service.AuthService,
	refreshService service.RefreshService,
	logger *slog.Logger,
) *AuthHandler {
	return &AuthHandler{
		authService:    authService,
		refreshService: refreshService,
		logger:         logger,
	}
}

func (h *AuthHandler) RegisterRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
	}
}

// -------------------- handlers --------------------

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid register body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Register(
		req.FirstName,
		req.LastName,
		req.Email,
		req.Password,
	)
	if err != nil {
		h.logger.Error("register failed", "email", req.Email, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.respondWithTokens(c, user.ID)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid login body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		h.logger.Warn("login failed", "email", req.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	h.respondWithTokens(c, user.ID)
}

// -------------------- helpers --------------------

func (h *AuthHandler) respondWithTokens(c *gin.Context, userID uint) {
	accessToken, err := jwt.GenerateAccessToken(userID)
	if err != nil {
		h.logger.Error("access token generation failed", "user_id", userID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "access token error"})
		return
	}

	refreshToken, err := h.refreshService.CreateForUser(userID)
	if err != nil {
		h.logger.Error("refresh token generation failed", "user_id", userID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "refresh token error"})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
