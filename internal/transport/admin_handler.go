package transport

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/IslamCHup/coworking-manager-project/internal/models"
	"github.com/IslamCHup/coworking-manager-project/internal/service"
	"github.com/IslamCHup/coworking-manager-project/internal/transport/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminHandler struct {
	userService    service.UserService
	bookingService service.BookingService
	logger         *slog.Logger
}

func NewAdminHandler(
	userService service.UserService,
	bookingService service.BookingService,
	logger *slog.Logger,
) *AdminHandler {
	return &AdminHandler{
		userService:    userService,
		bookingService: bookingService,
		logger:         logger,
	}
}

func (h *AdminHandler) RegisterRoutes(r *gin.Engine, adminService service.AdminService) {
	admin := r.Group("/admin", middleware.AdminBasicAuthMiddleware(adminService, h.logger))

	admin.GET("/login", h.Login)

	// Роуты для работы с пользователями
	admin.PUT("/users/:id", h.UpdateUser)
	admin.DELETE("/users/:id", h.DeleteUser)

	// Роуты для работы с букингами
	admin.PUT("/bookings/:id", h.UpdateBooking)
	admin.DELETE("/bookings/:id", h.DeleteBooking)
}

func (h *AdminHandler) Login(c *gin.Context) {
	adminID, exists := c.Get("admin_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "admin not found in context"})
		return
	}

	adminLogin, _ := c.Get("admin_login")
	c.JSON(http.StatusOK, gin.H{
		"id":    adminID,
		"login": adminLogin,
	})
}

func (h *AdminHandler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.logger.Warn("invalid user id", "id", idParam, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req models.UserUpdateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("UpdateUser invalid body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userService.UpdateUser(uint(userID), req); err != nil {
		h.logger.Error("UpdateUser failed", "user_id", userID, "error", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	h.logger.Info("UpdateUser success", "user_id", userID)
	c.JSON(http.StatusOK, gin.H{"message": "user updated successfully"})
}

func (h *AdminHandler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.logger.Warn("invalid user id", "id", idParam, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if err := h.userService.DeleteUser(uint(userID)); err != nil {
		h.logger.Error("DeleteUser failed", "user_id", userID, "error", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}

	h.logger.Info("DeleteUser success", "user_id", userID)
	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

func (h *AdminHandler) UpdateBooking(c *gin.Context) {
	idParam := c.Param("id")
	bookingID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.logger.Warn("invalid booking id", "id", idParam, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}

	var req models.BookingReqUpdateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("UpdateBooking invalid body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.bookingService.UpdateBook(uint(bookingID), &req); err != nil {
		h.logger.Error("UpdateBooking failed", "booking_id", bookingID, "error", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update booking"})
		return
	}

	h.logger.Info("UpdateBooking success", "booking_id", bookingID)
	c.JSON(http.StatusOK, gin.H{"message": "booking updated successfully"})
}

func (h *AdminHandler) DeleteBooking(c *gin.Context) {
	idParam := c.Param("id")
	bookingID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.logger.Warn("invalid booking id", "id", idParam, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid booking id"})
		return
	}

	if err := h.bookingService.DeleteBooking(uint(bookingID)); err != nil {
		h.logger.Error("DeleteBooking failed", "booking_id", bookingID, "error", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete booking"})
		return
	}

	h.logger.Info("DeleteBooking success", "booking_id", bookingID)
	c.JSON(http.StatusOK, gin.H{"message": "booking deleted successfully"})
}
