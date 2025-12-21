package repository

import "github.com/IslamCHup/coworking-manager-project/internal/models"

type ReviewRepository interface {
	CreateReview(req *models.Review) error
	UpdateReview(req *models.Review) error
	DeleteReview(id uint) error
	GetReview(id uint) (*models.Review, error)
}
