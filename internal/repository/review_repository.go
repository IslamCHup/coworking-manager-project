package repository

import (
	"errors"

	"github.com/IslamCHup/coworking-manager-project/internal/models"
	"gorm.io/gorm"
)

type ReviewRepository interface {
	CreateReview(req *models.Review) error
	// UpdateReview(req *models.Review) error
	// DeleteReview(id uint) error
	// GetReview(id uint) (*models.Review, error)
}

var ErrReviewNil = errors.New("review nil")

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db: db}
}
func (r *reviewRepository) GetByUserAndPlace(userID, placeID uint) (*models.Review, error) {
	var review models.Review
	err := r.db.Where("user_id = ? AND place_id = ?", userID, placeID).First(&review).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &review, nil
}
func (r *reviewRepository) CreateReview(review *models.Review) error {
	if review == nil {
		return ErrReviewNil
	}
	if err := r.db.Create(review).Error; err != nil {
		return err
	}
	return nil
}
