package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/IslamCHup/coworking-manager-project/internal/models"
	"github.com/IslamCHup/coworking-manager-project/internal/redis"
	"github.com/IslamCHup/coworking-manager-project/internal/repository"
	"gorm.io/gorm"
)

type BookingService interface {
	Create(id uint, req models.BookingReqDTO) (*models.Booking, error)
	GetBookingById(id uint) (*models.BookingResDTO, error)
	DeleteBooking(id uint) error
	ListBooking(filter *models.FilterBooking) ([]models.Booking, error)
	UpdateBook(id uint, req *models.BookingReqUpdateDTO) error
	UpdateStatus(id uint, status models.BookingStatusUpdateDTO) error
	UpdateBookingStatusWithBalance(id uint, newStatus models.BookingStatus) error
}

type bookingService struct {
	repo      repository.BookingRepository
	placeRepo repository.PlaceRepository
	db        *gorm.DB
	logger    *slog.Logger
	redis     *redis.Client
}

func NewBookingService(repo repository.BookingRepository, placeRepo repository.PlaceRepository, db *gorm.DB, logger *slog.Logger, redis *redis.Client) BookingService {
	return &bookingService{
		repo:      repo,
		placeRepo: placeRepo,
		db:        db,
		logger:    logger,
		redis:     redis,
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

	duration := end.Sub(start).Hours()
	if duration <= 0 {
		return nil, errors.New("неверный диапазон времени: время окончания должно быть позже времени начала")
	}

	if start.Weekday() == time.Saturday || start.Weekday() == time.Sunday {
		return nil, errors.New("нельзя бронировать на выходной день")
	}

	if start.Before(time.Now()) {
		return nil, errors.New("бронь просрочена")
	}

	if start.Hour() < 9 || end.Hour() > 18 {
		return nil, errors.New("мы работаем с 9 до 18 часов")
	}

	status := models.BookingActive
	filter := models.FilterBooking{
		PlaceID:   &req.PlaceID,
		StartTime: &start,
		EndTime:   &end,
		Status:    (*string)(&status),
	}

	bookings, err := s.repo.ListBooking(&filter)
	if err != nil {
		s.logger.Error("failed to list bookings for overlap check", "error", err)
		return nil, err
	}

	if len(bookings) > 0 {
		return nil, errors.New("это время занято другими")
	}

	booking := &models.Booking{
		UserID:    id,
		PlaceID:   req.PlaceID,
		StartTime: start,
		EndTime:   end,
		Status:    models.BookingNonActive,
	}

	place, err := s.placeRepo.GetPlaceByID(req.PlaceID)
	if err != nil {
		s.logger.Error("failed to get place for price calculation", "place_id", req.PlaceID, "error", err)
		return nil, errors.New("место не найдено")
	}

	booking.TotalPrice = int(duration * float64(place.PricePerHour))

	if err := s.repo.CreateBooking(booking); err != nil {
		s.logger.Error("Create booking failed", "error", err, "user_id", req.UserID, "place_id", req.PlaceID)
		return nil, err
	}

	s.logger.Info("Create booking success", "booking_id", booking.ID)

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
		return nil, errors.New("бронирование не найдено")
	}

	bookingResDTO := &models.BookingResDTO{
		UserID:     booking.UserID,
		PlaceID:    booking.PlaceID,
		StartTime:  booking.StartTime,
		EndTime:    booking.EndTime,
		TotalPrice: booking.TotalPrice,
		Status:     string(booking.Status),
		Place:      booking.Place,
	}

	bookingResDTO.User = &models.UserResponseDTO{
		ID:        booking.User.ID,
		Email:     booking.User.Email,
		FirstName: booking.User.FirstName,
		LastName:  booking.User.LastName,
		Balance:   booking.User.Balance,
	}
	s.logger.Info("get booking by id completed")

	return bookingResDTO, nil
}

func (s *bookingService) DeleteBooking(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		s.logger.Error("failed delete record")
		return err
	}

	s.logger.Info("booking deleted")

	return nil
}

func (s *bookingService) ListBooking(filter *models.FilterBooking) ([]models.Booking, error) {

	ctx := context.Background()

	cacheKey := buildBookingCacheKey(filter)

	if s.redis != nil {
		cached, err := s.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			s.logger.Info("ListBooking cache HIT", "key", cacheKey)

			var bookings []models.Booking
			if err := json.Unmarshal([]byte(cached), &bookings); err == nil {
				return bookings, nil
			}

			s.logger.Error(
				"failed to unmarshal bookings from cache",
				"key", cacheKey,
				"error", err,
			)
		} else {
			s.logger.Error(
				"redis GET error",
				"key", cacheKey,
				"error", err,
			)
		}
	}

	bookings, err := s.repo.ListBooking(filter)
	if err != nil {
		s.logger.Error("ListBooking repo error", "error", err)
		return nil, err
	}

	// 4. Кладём в Redis ТОЛЬКО если кешируемо
	if s.redis != nil {
		data, err := json.Marshal(bookings)
		if err != nil {
			s.logger.Error("failed to marshal bookings for cache", "error", err)
		} else {
			if err := s.redis.Set(ctx, cacheKey, data, 2*time.Minute).Err(); err != nil {
				s.logger.Error(
					"failed to set bookings cache",
					"key", cacheKey,
					"error", err,
				)
			} else {
				s.logger.Info("ListBooking cache SET", "key", cacheKey)
			}
		}
	}

	s.logger.Info("ListBooking success")
	return bookings, nil
}

func buildBookingCacheKey(filter *models.FilterBooking) string {
	parts := []string{"bookings", "v1"}

	if filter.PlaceID != nil {
		parts = append(parts, fmt.Sprintf("place:%d", *filter.PlaceID))
	}

	if filter.StartTime != nil && filter.EndTime != nil {
		parts = append(parts, fmt.Sprintf(
			"t:%d-%d", filter.StartTime.Unix(),filter.EndTime.Unix(),
		))
	}

	if filter.Status != nil {
		parts = append(parts, "status:"+strings.ToLower(*filter.Status))
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}

	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}

	if offset == 0 {
		parts = append(parts, fmt.Sprintf("lim:%d", limit))
	}

	return strings.Join(parts, ":")
}


func (s *bookingService) UpdateBook(id uint, req *models.BookingReqUpdateDTO) error {
	booking, err := s.repo.GetBookingById(id)
	if err != nil {
		s.logger.Error("failed to get booking for update", "id", id, "error", err)
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
			return errors.New("неверный диапазон времени: время окончания должно быть позже времени начала")
		}

		status := models.BookingActive
		filter := models.FilterBooking{
			PlaceID:   req.PlaceID,
			StartTime: &start,
			EndTime:   &end,
			Status:    (*string)(&status),
		}

		bookings, err := s.repo.ListBooking(&filter)
		if err != nil {
			s.logger.Error("failed to list bookings for overlap check", "error", err)
			return err
		}
		if len(bookings) > 0 {
			return errors.New("это время занято другими")
		}

		// Получаем информацию о месте для расчета цены
		place, err := s.placeRepo.GetPlaceByID(booking.PlaceID)
		if err != nil {
			s.logger.Error("failed to get place for price calculation", "place_id", booking.PlaceID, "error", err)
			return errors.New("место не найдено")
		}

		// Цена в копейках: часы * цена за час места в копейках
		booking.TotalPrice = int(durationHours * float64(place.PricePerHour))
	}

	if err := s.repo.UpdateBook(id, booking); err != nil {
		s.logger.Error("failed to update booking", "id", id, "error", err)
		return err
	}

	return nil
}

func (s *bookingService) UpdateStatus(id uint, status models.BookingStatusUpdateDTO) error {
	booking, err := s.repo.GetBookingById(id)
	if err != nil {
		return err
	}

	if status.Status == "" {
		return errors.New("статус не указан")
	}

	statusClear := models.BookingStatus(strings.ToLower(strings.TrimSpace(string(status.Status))))

	booking.Status = statusClear

	switch statusClear {
	case models.BookingActive, models.BookingNonActive, models.BookingCancelled:
		if err := s.repo.UpdateBook(booking.ID, booking); err != nil {
			return err
		}
		return nil
	default:
		return errors.New("неверный статус бронирования")
	}
}

func (s *bookingService) UpdateBookingStatusWithBalance(id uint, newStatus models.BookingStatus) error {
	// Начинаем транзакцию
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Получаем бронь в рамках транзакции
		var booking models.Booking
		if err := tx.Preload("User").Preload("Place").Where("id = ?", id).First(&booking).Error; err != nil {
			s.logger.Error("failed to get booking in transaction", "booking_id", id, "error", err)
			return err
		}

		oldStatus := booking.Status
		newStatusNormalized := models.BookingStatus(strings.ToLower(strings.TrimSpace(string(newStatus))))

		// Валидация статуса
		switch newStatusNormalized {
		case models.BookingActive, models.BookingNonActive, models.BookingCancelled:
			// OK
		default:
			return errors.New("неверный статус бронирования")
		}

		// Если статус не изменился, ничего не делаем
		if oldStatus == newStatusNormalized {
			s.logger.Info("booking status unchanged", "booking_id", id, "status", newStatusNormalized)
			return nil
		}

		// TotalPrice уже в копейках
		priceInCents := booking.TotalPrice

		// Логика для смены статуса на active
		if newStatusNormalized == models.BookingActive {
			// Проверяем баланс пользователя
			if booking.User.Balance < priceInCents {
				s.logger.Warn("insufficient balance", "user_id", booking.UserID, "balance", booking.User.Balance, "required", priceInCents)
				return errors.New("недостаточно средств")
			}

			// Списываем деньги с баланса пользователя
			if err := tx.Model(&models.User{}).Where("id = ?", booking.UserID).
				Update("balance", gorm.Expr("balance - ?", priceInCents)).Error; err != nil {
				s.logger.Error("failed to deduct balance", "user_id", booking.UserID, "error", err)
				return err
			}

			s.logger.Info("balance deducted", "user_id", booking.UserID, "amount", priceInCents)
		}

		// Логика для возврата денег при отмене активной брони
		if oldStatus == models.BookingActive && (newStatusNormalized == models.BookingCancelled || newStatusNormalized == models.BookingNonActive) {
			// Возвращаем деньги пользователю
			if err := tx.Model(&models.User{}).Where("id = ?", booking.UserID).
				Update("balance", gorm.Expr("balance + ?", priceInCents)).Error; err != nil {
				s.logger.Error("failed to refund balance", "user_id", booking.UserID, "error", err)
				return err
			}

			s.logger.Info("balance refunded", "user_id", booking.UserID, "amount", priceInCents)
		}

		// Обновляем статус брони
		booking.Status = newStatusNormalized
		if err := tx.Model(&models.Booking{}).Where("id = ?", id).Update("status", newStatusNormalized).Error; err != nil {
			s.logger.Error("failed to update booking status", "booking_id", id, "error", err)
			return err
		}

		s.logger.Info("booking status updated with balance transaction",
			"booking_id", id,
			"old_status", oldStatus,
			"new_status", newStatusNormalized,
			"user_id", booking.UserID)

		return nil
	})
}
