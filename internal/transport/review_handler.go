package transport

import (
	"log/slog"
	"net/http"

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
func (h ReviewHandler) RegisterRoutesReview(r *gin.Engine) {
	review := r.Group("/review")
	{
		review.POST("/", h.CreateReview)
		// booking.GET("/:id", h.GetByID)
		// booking.DELETE("/:id", h.DeleteBooking)
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
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "отзыв успешно создан",
		"review":  createdReview,
	})
}
