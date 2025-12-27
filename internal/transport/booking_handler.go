package transport

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/IslamCHup/coworking-manager-project/internal/middleware"
	"github.com/IslamCHup/coworking-manager-project/internal/models"
	"github.com/IslamCHup/coworking-manager-project/internal/service"
)

type BookingHandler struct {
	service service.BookingService
	logger  *slog.Logger
}

func NewBookingHandler(service service.BookingService, logger *slog.Logger) *BookingHandler {
	return &BookingHandler{service: service, logger: logger}
}

func (h *BookingHandler) RegisterRoutes(r *gin.Engine) {
	r.Use(middleware.JWTMiddleware())

	booking := r.Group("/bookings")
	{
		booking.GET("/:id", h.GetByID)
		booking.GET("/", h.ListBooking)
	}

	protected := r.Group("/bookings")
	protected.Use(middleware.RequireAuthMiddleware())
	{
		protected.POST("/", h.Create)
		protected.DELETE("/:id", h.DeleteBooking)
		protected.PATCH("/:id", h.Update)
		protected.PATCH("/status/:id", h.UpdateStatus)
	}
}

func (h *BookingHandler) GetByID(c *gin.Context) {
	userIDAny, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDAny.(uint)

	booking, err := h.service.GetBookingById(userID)
	if err != nil {
		h.logger.Error("GetBooking failed", "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("GetBooking success")

	c.JSON(http.StatusOK, booking)
}

func (h *BookingHandler) DeleteBooking(c *gin.Context) {
	id := c.MustGet("user_id").(uint)

	if err := h.service.DeleteBooking(uint(id)); err != nil {
		h.logger.Error("DeleteBooking failed", "error", err, "id", id)

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("DeleteBooking success", "booking_id", id)

	c.Status(http.StatusNoContent)
}

func (h *BookingHandler) Create(c *gin.Context) {

	var req models.BookingReqDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("CreateBooking invalid body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	booking, err := h.service.Create(c.MustGet("user_id").(uint), req)
	if err != nil {
		h.logger.Error("CreateBooking failed", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("CreateBooking success", "booking_id", booking.ID)

	c.JSON(http.StatusCreated, booking)
}

func (h *BookingHandler) ListBooking(c *gin.Context) {
	var q models.FilterBooking

	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filter := models.FilterBooking{
		Status:    q.Status,
		PriceMin:  q.PriceMin,
		PriceMax:  q.PriceMax,
		StartTime: q.StartTime,
		EndTime:   q.EndTime,
		Limit:     q.Limit,
		Offset:    q.Offset,
		SortBy:    q.SortBy,
		Order:     q.Order,
	}
	booking, err := h.service.ListBooking(&filter)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, booking)
}

func (h *BookingHandler) Update(c *gin.Context) {
	id := c.MustGet("user_id").(uint)

	var req models.BookingReqUpdateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("UpdateBooking invalid body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateBook(uint(id), &req); err != nil {
		h.logger.Error("UpdateBooking failed", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("UpdateBooking success", "booking_id", id)
	c.JSON(http.StatusOK, gin.H{"message": "обновлено"})
}

func (h *BookingHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		h.logger.Info("GetBooking invalid id param", "id", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "параметр id обязателен"})
		return
	}
	id, _ := strconv.ParseUint(idStr, 10, 64)

	var status models.BookingStatusUpdateDTO

	if err := c.ShouldBindJSON(&status); err != nil {
		h.logger.Error("UpdateBooking invalid body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateStatus(uint(id), status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "обновлено"})
}
