package repository

import (
	"log/slog"

	"github.com/IslamCHup/coworking-manager-project/internal/models"
	"gorm.io/gorm"
)

type BookingRepository interface {
	CreateBooking(req *models.Booking) error
	// UpdateBook(req *models.Booking) error
	Delete(id uint) error
	GetBookingById(id uint) (*models.Booking, error)
}

type bookingRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewBookingRepository(db *gorm.DB, logger *slog.Logger) BookingRepository {
	return &bookingRepository{db: db, logger: logger}
}

func (r *bookingRepository) CreateBooking(req *models.Booking) error {
	r.logger.Debug("creating booking", "user_id", req.UserID, "place_id", req.PlaceID, "start", req.StartTime, "end", req.EndTime)
	if err := r.db.Create(req).Error; err != nil {
		r.logger.Debug("failed to create booking", "error", err)
		r.logger.Error("CreateBooking failed", "error", err, "user_id", req.UserID, "place_id", req.PlaceID)
		return err
	}
	r.logger.Info("booking created", "user_id", req.UserID, "place_id", req.PlaceID)
	return nil
}

func (r *bookingRepository) GetBookingById(id uint) (*models.Booking, error) {
	var booking models.Booking
	r.logger.Debug("getting booking by id", "id", id)

	res := r.db.Preload("User").Preload("Place").Where("id = ?", id).First(&booking)

	if res.Error != nil {
		r.logger.Debug("delete failed", "error", res.Error)
		r.logger.Error("failed to delete record", "id", id, "error", res.Error)
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		r.logger.Info("no booking deleted", "id", id)
		return nil, gorm.ErrRecordNotFound
	}
	r.logger.Debug("booking successfully retrieved", "id", booking.ID)

	return &booking, nil
}

func (r *bookingRepository) Delete(id uint) error {
	r.logger.Debug("deleting booking", "id", id)

	res := r.db.Delete(&models.Booking{}, id)

	if res.RowsAffected == 0 {
		r.logger.Info("no booking deleted", "id", id)
		return gorm.ErrRecordNotFound
	}

	if res.Error != nil {
		r.logger.Debug("delete failed", "error", res.Error)
		r.logger.Error("failed to delete record", "id", id, "error", res.Error)
		return res.Error
	}
	r.logger.Info("booking deleted", "id", id, "rows_affected", res.RowsAffected)
	return nil
}
