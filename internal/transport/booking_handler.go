package transport

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

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

func (h BookingHandler) RegisterRoutes(r *gin.Engine) {
	booking := r.Group("/booking")
	{
		booking.POST("/", h.Create)
		booking.GET("/:id", h.GetByID)
		booking.DELETE("/:id", h.DeleteBooking)
		booking.GET("/", h.ListBooking)
		booking.PATCH("/:id", h.Update)
		booking.PATCH("/status/:id", h.UpdateStatus)
	}
}

func (h *BookingHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		h.logger.Info("GetBooking invalid id param", "id", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param is required"})
		return
	}

	id, _ := strconv.ParseUint(idStr, 10, 64)

	booking, err := h.service.GetBookingById(uint(id))
	if err != nil {
		h.logger.Error("GetBooking failed", "error", err, "id", id)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("GetBooking success")

	c.JSON(http.StatusOK, booking)
}

func (h *BookingHandler) DeleteBooking(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		h.logger.Info("DeleteBooking invalid id param", "id", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param is required"})
		return
	}

	id, _ := strconv.ParseUint(idStr, 10, 64)

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

	booking, err := h.service.Create(c.MustGet("UserID").(uint), req)
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
	idStr := c.Param("id")
	if idStr == "" {
		h.logger.Info("UpdateBooking invalid id param", "id", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param is required"})
		return
	}

	id, _ := strconv.ParseUint(idStr, 10, 64)

	var req models.BookingReqUpdateDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("UpdateBooking invalid body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateBooking(uint(id), &req); err != nil {
		h.logger.Error("UpdateBooking failed", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("UpdateBooking success", "booking_id", id)
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (h *BookingHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		h.logger.Info("GetBooking invalid id param", "id", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param is required"})
		return
	}
	id, _ := strconv.ParseUint(idStr, 10, 64)

	var status models.BookingStatusDTO

	if err := c.ShouldBindJSON(&status); err != nil {
		h.logger.Error("UpdateBooking invalid body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateStatus(uint(id), status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}
