package repository

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/IslamCHup/coworking-manager-project/internal/models"
	"gorm.io/gorm"
)

type BookingRepository interface {
	CreateBooking(req *models.Booking) error
	ListBooking(filter *models.FilterBooking) (*[]models.Booking, error)
	UpdateBook(id uint, req *models.Booking) error
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

func (r *bookingRepository) ListBooking(filter *models.FilterBooking) (*[]models.Booking, error) {
	if filter == nil {
		filter = &models.FilterBooking{
			Limit:  20,
			Offset: 0,
			SortBy: "start_time",
			Order:  "asc",
		}
	}
	if filter.Limit <= 0 || filter.Limit > 100 {
		filter.Limit = 20
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}
	if filter.SortBy == "" {
		filter.SortBy = "start_time"
	}
	if filter.Order == "" {
		filter.Order = "asc"
	}

	r.logger.Debug("ListBooking called", "filter", filter)

	var bookings *[]models.Booking

	r.logger.Debug("")

	query := r.db.Model(models.Booking{}).Preload("User").Preload("Place")

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.PriceMin != nil {
		query = query.Where("total_price >= ?", *filter.PriceMin)
	}
	if filter.PriceMax != nil {
		query = query.Where("total_price <= ?", *filter.PriceMax)
	}
	if filter.StartTime != nil && filter.EndTime != nil {
		query = query.Where("start_time >= ? AND end_time <= ?", *filter.StartTime, *filter.EndTime)
	} else if filter.StartTime != nil {
		query = query.Where("start_time >= ?", *filter.StartTime)
	} else if filter.EndTime != nil {
		query = query.Where("end_time <= ?", *filter.EndTime)
	}

	r.logger.Debug("filters applied", "status", filter.Status, "price_min", filter.PriceMin, "price_max", filter.PriceMax, "start", filter.StartTime, "end", filter.EndTime)

	allowed := map[string]bool{
		"start_time":  true,
		"total_price": true,
		"created_at":  true,
		"id":          true,
	}
	sortBy := strings.ToLower(filter.SortBy)
	if !allowed[sortBy] {
		sortBy = "start_time"
	}
	order := strings.ToLower(filter.Order)
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	query = query.Order(fmt.Sprintf("%s %s", sortBy, order))
	r.logger.Debug("query ready", "order", fmt.Sprintf("%s %s", sortBy, order), "limit", filter.Limit, "offset", filter.Offset)

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	query = query.Find(&bookings)

	if query.Error != nil {
		r.logger.Error("ListBooking failed", "err", query.Error)
		return nil, query.Error
	}
	r.logger.Info("ListBooking success", "count", len(*bookings), "limit", filter.Limit, "offset", filter.Offset)
	return bookings, nil
}

func (r *bookingRepository) UpdateBook(id uint, req *models.Booking) error{
	result := r.db.Model(models.Booking{}).Where("id = ?", id).Updates(req)
	if result.Error != nil {
		r.logger.Error("failed to update booking", "error", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.logger.Info("booking not found", "id", req.ID)
		return gorm.ErrRecordNotFound
	}
	return nil
}
