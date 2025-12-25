package transport

import (
	"log/slog"

	"github.com/IslamCHup/coworking-manager-project/internal/middleware"
	"github.com/IslamCHup/coworking-manager-project/internal/service"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.Engine,
	logger *slog.Logger,
	bookingService service.BookingService,
	adminService service.AdminService,
	userService service.UserService,
	authService service.AuthService,
	refreshService service.RefreshService,
	reviewService service.ReviewService,
) {
	bookingHandler := NewBookingHandler(bookingService, logger)
	bookingHandler.RegisterRoutes(router)
	authHandler := NewAuthHandler(authService, logger, refreshService)
	authHandler.RegisterRoutes(router)
	refreshHandler := NewRefreshHandler(refreshService, logger)
	refreshHandler.RegisterRoutes(router)
	userHandler := NewUserHandler(userService, logger)

	adminHandler := NewAdminHandler(userService, bookingService, logger)
	adminHandler.RegisterRoutes(router, adminService)

	reviewHandler := NewReviewHandler(reviewService, logger)

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())

	users := protected.Group("/users")
	userHandler.RegisterRoutes(users)

	reviews := protected.Group("/reviews")
	reviews.POST("/", reviewHandler.CreateReview)
}
