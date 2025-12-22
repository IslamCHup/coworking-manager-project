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
	}
}

func (h *BookingHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		h.logger.Info("GetBooking invalid id param", "id", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param is required"})
		return
	}

	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.Error("GetBooking invalid id format", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	booking, err := h.service.GetBookingById(uint(id64))
	if err != nil {
		h.logger.Error("GetBooking failed", "error", err, "id", id64)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("GetBooking success", "booking_id", booking.ID)

	c.JSON(http.StatusOK, booking)
}

func (h *BookingHandler) DeleteBooking(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		h.logger.Info("DeleteBooking invalid id param", "id", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "id param is required"})
		return
	}

	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.Error("DeleteBooking invalid id format", "id", idStr, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DeleteBooking(uint(id64)); err != nil {
		h.logger.Error("DeleteBooking failed", "error", err, "id", id64)

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("DeleteBooking success", "booking_id", id64)

	c.Status(http.StatusNoContent)
}

func (h *BookingHandler) Create(c *gin.Context) {

	var req models.BookingReqDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		if h.logger != nil {
			h.logger.Error("CreateBooking invalid body", "error", err)
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	booking, err := h.service.Create(req)
	if err != nil {
		if h.logger != nil {
			h.logger.Error("CreateBooking failed", "error", err)
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if h.logger != nil {
		h.logger.Info("CreateBooking success", "booking_id", booking.ID)
	}

	c.JSON(http.StatusCreated, booking)
}
