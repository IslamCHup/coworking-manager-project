package service

import (
	"errors"

	"github.com/IslamCHup/coworking-manager-project/internal/models"
)
var (
	ErrUserNotFound = errors.New("Пользователь не оставил отзыв ")
)

type ReviewService interface {
	CreateReview(req *models.Review) (*models.Review , error)
	GetReviewId(id uint) (*models.Review,error)
	UpdateReview(id uint, req models.Review) (*models.Review,error)
	DeleteReview(id uint)error
}