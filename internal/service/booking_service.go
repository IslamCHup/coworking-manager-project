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
	GetBookingById(id uint) (*models.BookingResDTO, error)
	UpdateBooking(id uint, req models.BookingUpdateDTO) error
	DeleteBooking(id uint) error
	ListBooking(filter *models.FilterBooking) (*[]models.Booking, error)
}

type bookingService struct {
	repo      repository.BookingRepository
	placeRepo repository.PlaceRepository
	logger    *slog.Logger
}

func NewBookingService(repo repository.BookingRepository, placeRepo repository.PlaceRepository, logger *slog.Logger) BookingService {
	return &bookingService{
		repo:      repo,
		placeRepo: placeRepo,
		logger:    logger,
	}
}

func (s *bookingService) Create(req models.BookingReqDTO) (*models.Booking, error) {
	// Загружаем Place для получения цены за час
	place, err := s.placeRepo.GetPlaceByID(req.PlaceID)
	if err != nil {
		s.logger.Error("Place not found", "place_id", req.PlaceID, "error", err)
		return nil, errors.New("place not found")
	}

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

	// Используем PricePerHour из Place
	price := durationHours * place.PricePerHour
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

func (s *bookingService) GetBookingById(id uint) (*models.BookingResDTO, error) {
	booking, err := s.repo.GetBookingById(id)

	if err != nil {
		s.logger.Error("failed to get booking from repository")
		return nil, err
	}

	if booking == nil {
		s.logger.Error("booking not found")
		return nil, errors.New("booking not found")
	}

	bookingResDTO := &models.BookingResDTO{
		ID:         booking.ID,
		UserID:     booking.UserID,
		PlaceID:    booking.PlaceID,
		StartTime:  booking.StartTime,
		EndTime:    booking.EndTime,
		TotalPrice: booking.TotalPrice,
		Status:     string(booking.Status),
		//place сделать дто после как ее доделают
		Place: booking.Place,
	}

	bookingResDTO.User = &models.UserResponseDTO{
		ID:        booking.User.ID,
		Phone:     booking.User.Phone,
		FirstName: booking.User.FirstName,
		LastName:  booking.User.LastName,
	}
	s.logger.Info("get booking by id completed")

	return bookingResDTO, nil
}

func (s *bookingService) UpdateBooking(id uint, req models.BookingUpdateDTO) error {
	booking, err := s.repo.GetBookingById(id)
	if err != nil {
		s.logger.Error(
			"booking not found",
			"booking_id", id,
			"error", err,
		)
		return err
	}

	if req.UserID != nil {
		booking.UserID = *req.UserID
	}

	// Если изменился PlaceID, загружаем новый Place
	if req.PlaceID != nil {
		booking.PlaceID = *req.PlaceID
		// Загружаем новый Place для получения актуальной цены
		newPlace, err := s.placeRepo.GetPlaceByID(*req.PlaceID)
		if err != nil {
			s.logger.Error(
				"place not found",
				"place_id", *req.PlaceID,
				"error", err,
			)
			return errors.New("place not found")
		}
		booking.Place = newPlace
	}

	if req.StartTime != nil {
		booking.StartTime = *req.StartTime
	}

	if req.EndTime != nil {
		booking.EndTime = *req.EndTime
	}

	if req.Status != nil {
		booking.Status = *req.Status
	}

	// Пересчитываем цену, если изменилось время или PlaceID
	if req.StartTime != nil || req.EndTime != nil || req.PlaceID != nil {
		durationHours := booking.EndTime.Sub(booking.StartTime).Hours()
		if durationHours <= 0 {
			return errors.New("invalid booking time range: end must be after start")
		}

		// Используем PricePerHour из Place
		// Если Place не загружен (не менялся PlaceID), загружаем его
		if booking.Place == nil {
			place, err := s.placeRepo.GetPlaceByID(booking.PlaceID)
			if err != nil {
				s.logger.Error(
					"place not found",
					"place_id", booking.PlaceID,
					"error", err,
				)
				return errors.New("place not found")
			}
			booking.Place = place
		}

		pricePerHour := booking.Place.PricePerHour
		price := durationHours * pricePerHour
		booking.TotalPrice = math.Round(price*100) / 100
	}

	if err := s.repo.UpdateBooking(booking); err != nil {
		s.logger.Error(
			"UpdateBooking failed",
			"booking_id", id,
			"error", err,
		)
		return err
	}

	s.logger.Info("UpdateBooking success", "booking_id", id)
	return nil
}

func (s *bookingService) DeleteBooking(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		s.logger.Error("failed delete record")
		return err
	}

	s.logger.Info("booking deleted")

	return nil
}

func (s *bookingService) ListBooking(filter *models.FilterBooking) (*[]models.Booking, error) {
	bookings, err := s.repo.ListBooking(filter)

	if err != nil {
		s.logger.Error("")
		return nil, err
	}

	s.logger.Info("ListBooking success")
	return bookings, nil
}
