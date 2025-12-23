package transport

import (
	"log/slog"
	"net/http"

	"github.com/IslamCHup/coworking-manager-project/internal/models"
	"github.com/IslamCHup/coworking-manager-project/internal/service"
	"github.com/gin-gonic/gin"
)

type RefreshRequestDTO struct {
	RefreshToken string
}

type RefreshHandler struct {
	authService    service.AuthService
	refreshService service.RefreshService
	logger         *slog.Logger
}

func NewRefreshHandler(
	authService service.AuthService,
	refreshService service.RefreshService,
	logger *slog.Logger,
) *RefreshHandler {
	return &RefreshHandler{
		authService:    authService,
		refreshService: refreshService,
		logger:         logger,
	}
}

func (h *RefreshHandler) Refresh(c *gin.Context) {
	var req models.RefreshRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Refresh invalid body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	access, refresh, err := h.refreshService.Refresh(req.RefreshToken)
	if err != nil {
		h.logger.Warn("Refresh failed", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Refresh success")
	c.JSON(http.StatusOK, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func (h *RefreshHandler) Logout(c *gin.Context) {
	var req models.RefreshRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Logout invalid body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.refreshService.Logout(req.RefreshToken); err != nil {
		h.logger.Error("Logout failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "logout failed"})
		return
	}

	h.logger.Info("Logout success")
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}
