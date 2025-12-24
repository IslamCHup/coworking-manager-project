package service

import (
	"errors"
	"log/slog"
	"math"
	"strings"
	"time"

	"github.com/IslamCHup/coworking-manager-project/internal/models"
	"github.com/IslamCHup/coworking-manager-project/internal/repository"
)

type BookingService interface {
	Create(id uint, req models.BookingReqDTO) (*models.Booking, error)
	GetBookingById(id uint) (*models.BookingResDTO, error)
	UpdateBooking(id uint, req models.BookingUpdateDTO) error
	DeleteBooking(id uint) error
	ListBooking(filter *models.FilterBooking) (*[]models.Booking, error)
	UpdateBooking(id uint, req *models.BookingReqUpdateDTO) error
	UpdateStatus(id uint, status models.BookingStatusDTO) error 
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

func (s *bookingService) Create(id uint, req models.BookingReqDTO) (*models.Booking, error) {
	start, err := time.Parse("2006-01-02 15", req.StartTime)
	if err != nil {
		return nil, errors.New("неправильный формат времени, нужен YYYY-MM-DD HH")
	}

	end, err := time.Parse("2006-01-02 15", req.EndTime)
	if err != nil {
		return nil, errors.New("неправильный формат времени, нужен YYYY-MM-DD HH")
	}

	duration := end.Sub(start)
	if duration <= 0 {
		return nil, errors.New("invalid booking time range: end must be after start")
	}

	if start.Weekday() == time.Saturday || start.Weekday() == time.Sunday {
		return nil, errors.New("нельзя бронировать на выходной день")
	}

	if start.Before(time.Now()) {
		return nil, errors.New("бронь просрочена")
	}

	if start.Hour() < 9 || end.Hour() > 17 {
		return nil, errors.New("мы работаем с 9 до 18 часов")
	}

	bookings, err := s.repo.ListBooking(nil)
	if err != nil {
		if s.logger != nil {
			s.logger.Error("failed to list bookings for overlap check", "error", err)
		}
		return nil, err
	}
	if bookings != nil {
		for _, v := range *bookings {
			if v.PlaceID != req.PlaceID {
				continue
			}
			existingStart := v.StartTime
			existingEnd := v.EndTime
			newStart := start
			newEnd := end
			if existingStart.Before(newEnd) && existingEnd.After(newStart) && v.Status == models.BookingActive {
				return nil, errors.New("это время занято другими")
			}
		}
	}

	booking := &models.Booking{
		UserID:    id,
		PlaceID:   req.PlaceID,
		StartTime: start,
		EndTime:   end,
		Status:    models.BookingNonActive,
	}

	durationHours := booking.EndTime.Sub(booking.StartTime).Hours()
	price := durationHours * float64(models.PriceHour)

	booking.TotalPrice = math.Round(price*100) / 100

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

func (s *bookingService) UpdateBooking(id uint, req *models.BookingReqUpdateDTO) error {
	booking, err := s.repo.GetBookingById(id)
	if err != nil {
		if s.logger != nil {
			s.logger.Error("failed to get booking for update", "id", id, "error", err)
		}
		return err
	}

	start, err := time.Parse("2006-01-02 15", *req.StartTime)

	if err != nil {
		return errors.New("неправильный формат времени, нужен YYYY-MM-DD HH")
	}

	end, err := time.Parse("2006-01-02 15", *req.EndTime)
	if err != nil {
		return errors.New("неправильный формат времени, нужен YYYY-MM-DD HH")
	}

	if req.UserID != nil {
		booking.UserID = *req.UserID
	}
	if req.PlaceID != nil {
		booking.PlaceID = *req.PlaceID
	}
	if req.StartTime != nil {
		booking.StartTime = start
	}
	if req.EndTime != nil {
		booking.EndTime = end
	}

	if req.StartTime != nil || req.EndTime != nil {
		durationHours := booking.EndTime.Sub(booking.StartTime).Hours()
		if durationHours <= 0 {
			return errors.New("invalid booking time range: end must be after start")
		}

		bookings, err := s.repo.ListBooking(nil)
		if err != nil {
			if s.logger != nil {
				s.logger.Error("failed to list bookings for overlap check", "error", err)
			}
			return err
		}
		if bookings != nil {
			for _, v := range *bookings {
				if v.ID == booking.ID || v.PlaceID != booking.PlaceID {
					continue
				}
				if v.StartTime.Before(booking.EndTime) && v.EndTime.After(booking.StartTime) && v.Status == models.BookingActive {
					return errors.New("это время занято другими")
				}
			}
		}

		price := durationHours * float64(models.PriceHour)
		booking.TotalPrice = math.Round(price*100) / 100
	}

	if err := s.repo.UpdateBook(id, booking); err != nil {
		if s.logger != nil {
			s.logger.Error("failed to update booking", "id", id, "error", err)
		}
		return err
	}

	return nil
}

func (s *bookingService) UpdateStatus(id uint, status models.BookingStatusDTO) error {
	booking, err := s.repo.GetBookingById(id)
	if err != nil {
		return err
	}

	if status.Status == "" {
		return errors.New("empty status")
	}

	statusClear := models.BookingStatus(strings.ToLower(strings.TrimSpace(string(status.Status))))
	
	booking.Status = statusClear
	
	switch statusClear {
	case models.BookingActive, models.BookingNonActive, models.BookingCancelled:
		if err := s.repo.UpdateBook(booking.ID, booking); err != nil{
			return err
		}
		return nil
	default:
		return errors.New("invalid booking status")
	}
}
