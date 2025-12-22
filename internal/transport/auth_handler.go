package transport

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/IslamCHup/coworking-manager-project/internal/models"
	"github.com/IslamCHup/coworking-manager-project/internal/service"
)

type AuthHandler struct {
	service service.AuthService
	logger  *slog.Logger
}

func NewAuthHandler(
	service service.AuthService,
	logger *slog.Logger,
) *AuthHandler {
	return &AuthHandler{
		service: service,
		logger:  logger,
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
	var req models.PhoneVerifyDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("VerifyPhoneCode invalid body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.service.VerifyPhoneCode(req.Phone, req.Code)
	if err != nil {
		h.logger.Warn(
			"VerifyPhoneCode failed",
			"phone", req.Phone,
			"error", err,
		)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info(
		"VerifyPhoneCode success",
		"user_id", userID,
		"phone", req.Phone,
	)

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
	})
}
