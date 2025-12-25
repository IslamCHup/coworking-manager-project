package transport

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

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

	admin.GET("/users", h.GetAllUsers)
	admin.PUT("/users/:id", h.UpdateUser)
	admin.DELETE("/users/:id", h.DeleteUser)
	admin.PATCH("/users/:id/balance", h.UpdateUserBalance)

	admin.PUT("/bookings/:id", h.UpdateBooking)
	admin.DELETE("/bookings/:id", h.DeleteBooking)

	admin.PUT("/status/booking/:id", h.AdminUpdateBookingStatus)
}

func (h *AdminHandler) Login(c *gin.Context) {
	adminID, exists := c.Get("admin_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "администратор не найден в контексте"})
		return
	}

	adminLogin, _ := c.Get("admin_login")
	c.JSON(http.StatusOK, gin.H{
		"id":    adminID,
		"login": adminLogin,
	})
}

func (h *AdminHandler) GetAllUsers(c *gin.Context) {
    users, err := h.userService.GetAllUsers()
    if err != nil {
        h.logger.Error("Admin GetAllUsers failed", "error", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось получить пользователей"})
        return
    }
    c.JSON(http.StatusOK, users)
}

func (h *AdminHandler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.logger.Warn("invalid user id", "id", idParam, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID пользователя"})
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
			c.JSON(http.StatusNotFound, gin.H{"error": "пользователь не найден"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось обновить пользователя"})
		return
	}

	h.logger.Info("UpdateUser success", "user_id", userID)
	c.JSON(http.StatusOK, gin.H{"message": "пользователь успешно обновлен"})
}

func (h *AdminHandler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.logger.Warn("invalid user id", "id", idParam, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID пользователя"})
		return
	}

	if err := h.userService.DeleteUser(uint(userID)); err != nil {
		h.logger.Error("DeleteUser failed", "user_id", userID, "error", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "пользователь не найден"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось удалить пользователя"})
		return
	}

	h.logger.Info("DeleteUser success", "user_id", userID)
	c.JSON(http.StatusOK, gin.H{"message": "пользователь успешно удален"})
}

func (h *AdminHandler) UpdateBooking(c *gin.Context) {
	idParam := c.Param("id")
	bookingID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.logger.Warn("invalid booking id", "id", idParam, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID бронирования"})
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
			c.JSON(http.StatusNotFound, gin.H{"error": "бронирование не найдено"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось обновить бронирование"})
		return
	}

	h.logger.Info("UpdateBooking success", "booking_id", bookingID)
	c.JSON(http.StatusOK, gin.H{"message": "бронирование успешно обновлено"})
}

func (h *AdminHandler) DeleteBooking(c *gin.Context) {
	idParam := c.Param("id")
	bookingID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.logger.Warn("invalid booking id", "id", idParam, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID бронирования"})
		return
	}

	if err := h.bookingService.DeleteBooking(uint(bookingID)); err != nil {
		h.logger.Error("DeleteBooking failed", "booking_id", bookingID, "error", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "бронирование не найдено"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось удалить бронирование"})
		return
	}

	h.logger.Info("DeleteBooking success", "booking_id", bookingID)
	c.JSON(http.StatusOK, gin.H{"message": "бронирование успешно удалено"})
}

func (h *AdminHandler) AdminUpdateBookingStatus(c *gin.Context) {
	idParam := c.Param("id")
	bookingID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.logger.Warn("invalid booking id", "id", idParam, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID бронирования"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("AdminUpdateBookingStatus invalid body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "неверное тело запроса",
			"details": err.Error(),
		})
		return
	}

	// Валидируем статус вручную
	statusStr := strings.ToLower(strings.TrimSpace(req.Status))
	var bookingStatus models.BookingStatus
	switch statusStr {
	case "active":
		bookingStatus = models.BookingActive
	case "non_active":
		bookingStatus = models.BookingNonActive
	case "cancelled":
		bookingStatus = models.BookingCancelled
	default:
		h.logger.Warn("AdminUpdateBookingStatus invalid status value", "status", req.Status)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "неверное значение статуса",
			"details": "статус должен быть одним из: active, non_active, cancelled",
		})
		return
	}

	// Получаем информацию о букинге для детального сообщения об ошибке
	bookingInfo, _ := h.bookingService.GetBookingById(uint(bookingID))

	if err := h.bookingService.UpdateBookingStatusWithBalance(uint(bookingID), bookingStatus); err != nil {
		h.logger.Error("AdminUpdateBookingStatus failed", "booking_id", bookingID, "error", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "бронирование не найдено"})
			return
		}
		if err.Error() == "insufficient balance" {
			// Формируем детальное сообщение об ошибке
			errorResponse := gin.H{
				"error":   "недостаточно средств",
				"message": "У пользователя недостаточно средств для активации этого бронирования",
			}

			// Добавляем детали, если удалось получить информацию о букинге
			if bookingInfo != nil && bookingInfo.User != nil {
				errorResponse["details"] = gin.H{
					"user_id":             bookingInfo.UserID,
					"user_balance":        bookingInfo.User.Balance,
					"required":            bookingInfo.TotalPrice,
					"shortage":            bookingInfo.TotalPrice - bookingInfo.User.Balance,
					"user_balance_rubles": float64(bookingInfo.User.Balance) / 100,
					"required_rubles":     float64(bookingInfo.TotalPrice) / 100,
					"shortage_rubles":     float64(bookingInfo.TotalPrice-bookingInfo.User.Balance) / 100,
				}
			}

			c.JSON(http.StatusBadRequest, errorResponse)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось обновить статус бронирования"})
		return
	}

	h.logger.Info("AdminUpdateBookingStatus success", "booking_id", bookingID, "status", bookingStatus)
	c.JSON(http.StatusOK, gin.H{"message": "статус бронирования успешно обновлен"})
}

func (h *AdminHandler) UpdateUserBalance(c *gin.Context) {
    idParam := c.Param("id")
    userID64, err := strconv.ParseUint(idParam, 10, 32)
    if err != nil {
        h.logger.Warn("invalid user id", "id", idParam, "error", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID пользователя"})
        return
    }

    var req models.UpdateBalanceDTO
    if err := c.ShouldBindJSON(&req); err != nil {
        h.logger.Warn("UpdateUserBalance invalid body", "error", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }


    if req.Amount == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "amount не может быть 0"})
        return
    }

    if err := h.userService.UpdateUserBalance(uint(userID64), req.Amount); err != nil {
        h.logger.Error("UpdateUserBalance failed", "user_id", userID64, "amount", req.Amount, "error", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось изменить баланс"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "баланс обновлен"})
}