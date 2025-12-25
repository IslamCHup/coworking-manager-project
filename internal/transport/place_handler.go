package transport

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/IslamCHup/coworking-manager-project/internal/models"
	"github.com/IslamCHup/coworking-manager-project/internal/service"
)

type PlaceHandler struct {
	service service.PlaceService
	logger  *slog.Logger
}

func NewPlaceHandler(s service.PlaceService, logger *slog.Logger) *PlaceHandler {
	return &PlaceHandler{service: s, logger: logger}
}

func (h *PlaceHandler) RegisterRoutes(r *gin.Engine) {
	places := r.Group("/places")
	{
		places.GET("/", h.ListPlaces)
		places.GET("/free", h.ListFreePlaces)
		places.GET(":id", h.GetByID)
	}
}

func (h *PlaceHandler) ListPlaces(c *gin.Context) {
	var q models.FilterPlace
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	places, err := h.service.ListPlaces(&q)
	if err != nil {
		h.logger.Error("ListPlaces failed", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, places)
}

func (h *PlaceHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id обязателен"})
		return
	}
	id64, _ := strconv.ParseUint(idStr, 10, 64)

	place, err := h.service.GetPlaceByID(uint(id64))
	if err != nil {
		h.logger.Error("GetPlaceByID failed", "error", err, "id", id64)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, place)
}

func (h *PlaceHandler) ListFreePlaces(c *gin.Context) {
	var q models.FilterPlace
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	places, err := h.service.ListFreePlaces(&q)
	if err != nil {
		h.logger.Error("ListFreePlaces failed", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, places)
}