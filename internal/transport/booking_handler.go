package transport

import (
	"log/slog"
	"net/http"

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

func (h BookingHandler) RegisterRoutes(r *gin.Engine){
	booking := r.Group("/booking")
	{
		booking.POST("/", h.Create)
	}
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
