package transport

import (
	"log/slog"

	"github.com/IslamCHup/coworking-manager-project/internal/service"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.Engine,
	logger *slog.Logger,
	bookingService service.BookingService,
	adminService service.AdminService,
	userService service.UserService,
) {
	bookingHandler := NewBookingHandler(bookingService, logger)
	bookingHandler.RegisterRoutes(router)

	adminHandler := NewAdminHandler(userService, bookingService, logger)
	adminHandler.RegisterRoutes(router, adminService)
}
