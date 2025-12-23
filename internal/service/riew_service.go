package service

import (
	"errors"

	"time"

	"github.com/IslamCHup/coworking-manager-project/internal/models"
	"github.com/IslamCHup/coworking-manager-project/internal/repository"
	"gorm.io/gorm"
)

type ReviewService interface {
	CreateReview(req *models.Review) (*models.Review, error)
	// GetReviewId(id uint) (*models.Review,error)
	// UpdateReview(id uint, req models.Review) (*models.Review,error)
	// DeleteReview(id uint)error
}
type reviewService struct {
	db     *gorm.DB
	review repository.ReviewRepository
}

func NewReviewService(db *gorm.DB, review repository.ReviewRepository) ReviewService {
	return &reviewService{db: db, review: review}
}
func (s *reviewService) CreateReview(req *models.Review) (*models.Review, error) {
	if req.Rating < 1 || req.Rating > 5 {
		return nil, errors.New("рейтинг должен быть от 1 до 5 ")
	}
	if req.Text == "" {
		return nil, errors.New("текст не может быть пустым")
	}
	if len(req.Text) < 5 {
		return nil, errors.New("текст не может быть меньше 5 симловов ")
	}
	if len(req.Text) > 300 {
		return nil, errors.New("текст не может быть больше 300 симловов ")
	}
	if req.UserID == 0 {
		return nil, errors.New("пользователь не найден")
	}
	if req.PlaceID == 0 {
		return nil, errors.New("место не найдено")
	}
	review := &models.Review{
		UserID:    req.UserID,
		PlaceID:   req.PlaceID,
		Rating:    req.Rating,
		Text:      req.Text,
		CreatedAt: time.Now(),
	}
	err := s.review.CreateReview(review)
	if err != nil {
		return nil, err
	}
	return review, nil

}
