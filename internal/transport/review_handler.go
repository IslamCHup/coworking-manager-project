package transport

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/IslamCHup/coworking-manager-project/internal/models"
	"github.com/IslamCHup/coworking-manager-project/internal/service"
	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	review service.ReviewService
	logger *slog.Logger
}

func NewReviewHandler(review service.ReviewService, logger *slog.Logger) *ReviewHandler {
	return &ReviewHandler{review: review, logger: logger}
}

func (h *ReviewHandler) RegisterRoutesReview(r *gin.Engine) {
	review := r.Group("/review")
	{
		review.POST("/", h.CreateReview)
		review.GET("/:id", h.GetByID)
		 review.PUT("/:id", h.UpdateReview)
	}
}

func (h *ReviewHandler) CreateReview(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "пользователь не авторизован",
		})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ошибка сервера: неверный тип user_id",
		})
		return
	}

	if userIDUint == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "неверный идентификатор пользователя",
		})
		return
	}

	var req models.Review
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "неверный формат данных: " + err.Error(),
		})
		return
	}

	req.UserID = userIDUint

	createdReview, err := h.review.CreateReview(&req)
	if err != nil {
		h.logger.Error("failed to create review", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "не удалось создать отзыв",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "отзыв успешно создан",
		"review":  createdReview,
	})
}

func (h *ReviewHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	review, err := h.review.GetReviewId(uint(id))
	if err != nil {
		h.logger.Error("failed to get review by id", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "отзыв не найден"})
		return
	}

	c.JSON(http.StatusOK, review)
}

func (h *ReviewHandler) UpdateReview(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "пользователь не авторизован",
		})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok || userIDUint == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "неверный идентификатор пользователя",
		})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req models.UpdateReviewDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review, err := h.review.GetReviewId(uint(id))
	if err != nil {
		h.logger.Error("failed to get review", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "отзыв не найден"})
		return
	}

	if review.UserID != userIDUint {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "нет доступа к этому отзыву",
		})
		return
	}

	updatedReview, err := h.review.UpdateReview(uint(id), req)
	if err != nil {
		h.logger.Error("failed to update review", "id", id, "user_id", userIDUint, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось обновить отзыв"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "отзыв успешно обновлен",
		"review":  updatedReview,
	})
}