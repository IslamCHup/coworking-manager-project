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
	service        service.AuthService
	logger         *slog.Logger
	refreshService service.RefreshService
}

func NewAuthHandler(
	service service.AuthService,
	logger *slog.Logger,
	refreshService service.RefreshService,
) *AuthHandler {
	return &AuthHandler{
		service:        service,
		logger:         logger,
		refreshService: refreshService,
	}
}

func (h *AuthHandler) RegisterRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/phone/request", h.RequestPhoneCode)
		auth.POST("/phone/verify", h.VerifyPhoneCode)
	}
}

func (h *AuthHandler) RequestPhoneCode(c *gin.Context) {
	var req models.PhoneRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("RequestPhoneCode invalid body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.RequestPhoneCode(req.Phone); err != nil {
		h.logger.Error(
			"RequestPhoneCode failed",
			"phone", req.Phone,
			"error", err,
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("RequestPhoneCode success", "phone", req.Phone)
	c.JSON(http.StatusOK, gin.H{"message": "code sent"})
}

func (h *AuthHandler) VerifyPhoneCode(c *gin.Context) {
	var dto models.PhoneVerifyDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.service.VerifyPhoneCode(dto.Phone, dto.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := jwt.GenerateAccessToken(userID)
	if err != nil {
		h.logger.Error("GenerateAccessToken failed", "user_id", userID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		return
	}

	refreshToken, err := h.refreshService.CreateForUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
