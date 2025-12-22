package service

import (
	"errors"
	"log/slog"
	"math"

	"github.com/IslamCHup/coworking-manager-project/internal/models"
	"github.com/IslamCHup/coworking-manager-project/internal/repository"
)

type BookingService interface {
	Create(req models.BookingReqDTO) (*models.Booking, error)
}

type bookingService struct {
	repo   repository.BookingRepository
	logger *slog.Logger
}

func NewBookingService(repo repository.BookingRepository, logger *slog.Logger) BookingService {
	return &bookingService{repo: repo, logger: logger}
}

func (s *bookingService) Create(req models.BookingReqDTO) (*models.Booking, error) {
// исправить если другой прайс
	priceHour := 100

	booking := &models.Booking{
		UserID:     req.UserID,
		PlaceID:    req.PlaceID,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		TotalPrice: req.TotalPrice,
		Status:     req.Status,
	}

	durationHours := booking.EndTime.Sub(booking.StartTime).Hours()
	if durationHours <= 0 {
		return nil, errors.New("invalid booking time range: end must be after start")
	}

	price := durationHours * float64(priceHour)
	
	booking.TotalPrice = math.Round(price*100) / 100

	if booking.Status == "" {
		booking.Status = models.BookingActive
	}

	if err := s.repo.CreateBooking(booking); err != nil {
		if s.logger != nil {
			s.logger.Error("Create booking failed", "error", err, "user_id", req.UserID, "place_id", req.PlaceID)
		}
		return nil, err
	}

	if s.logger != nil {
		s.logger.Info("Create booking success", "booking_id", booking.ID)
	}

	return booking, nil
}
