package repository

import (
	"log/slog"

	"github.com/IslamCHup/coworking-manager-project/internal/models"
	"gorm.io/gorm"
)

type BookingRepository interface {
	CreateBooking(req *models.Booking) error
	// UpdateBook(req *models.Booking) error
	// DeleteBook(id uint) error
	// GetBookId(id uint) (*models.Booking, error)
}

type bookingRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewBookingRepository(db *gorm.DB, logger *slog.Logger) BookingRepository {
	return &bookingRepository{db: db, logger: logger}
}

func (r *bookingRepository) CreateBooking(req *models.Booking) error {
	if err := r.db.Create(req).Error; err != nil {
		if r.logger != nil {
			r.logger.Error("failed to create booking", "error", err)
		}
		return err
	}
	return nil
}
