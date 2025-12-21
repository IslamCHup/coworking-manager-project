package repository

import "github.com/IslamCHup/coworking-manager-project/internal/models"

type ReviewRepository interface {
	CreateReview(review *models.Review) error
	UpdateReview(review *models.Review)
	DeleteReview(id uint) error
	GetReview(id uint) (*models.Review, error)
}
