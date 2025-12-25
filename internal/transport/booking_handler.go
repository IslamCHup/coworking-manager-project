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

func (h BookingHandler) RegisterRoutes(r *gin.RouterGroup) {
		r.POST("/", h.Create)
		r.GET("/:id", h.GetByID)
		r.DELETE("/:id", h.DeleteBooking)
		r.PATCH("/:id", h.Update)
	}


func (h *BookingHandler) GetByID(c *gin.Context) {
	id := c.MustGet("userID").(uint)

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
	id := c.MustGet("userID").(uint)

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

func (h *BookingHandler) Update(c *gin.Context) {
	id := c.MustGet("userID").(uint)


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

