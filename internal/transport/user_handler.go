package transport

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/IslamCHup/coworking-manager-project/internal/middleware"
	"github.com/IslamCHup/coworking-manager-project/internal/models"
	"github.com/IslamCHup/coworking-manager-project/internal/service"
)

type UserHandler struct {
	service service.UserService
	logger  *slog.Logger
}

func NewUserHandler(
	service service.UserService,
	logger *slog.Logger,
) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

func (h *UserHandler) RegisterRoutes(r *gin.RouterGroup) {
	protected := r.Group("/")
	protected.Use(middleware.RequireAuthMiddleware())

	protected.GET("/me", h.GetUser)
	protected.PATCH("/me", h.UpdateUser)
	// r.GET("/", h.GetAllUsers) // только для admin
}

func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	user, err := h.service.GetUserByID(userID)
	if err != nil {
		h.logger.Error("GetUser failed", "user_id", userID, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "пользователь не найден"})
		return
	}

	h.logger.Info("GetUser success", "user_id", userID)
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var req models.UserUpdateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("UpdateUser invalid body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateUser(userID, req); err != nil {
		h.logger.Error("UpdateUser failed", "user_id", userID, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("UpdateUser success", "user_id", userID)
	c.JSON(http.StatusOK, gin.H{"message": "пользователь обновлен"})
}
