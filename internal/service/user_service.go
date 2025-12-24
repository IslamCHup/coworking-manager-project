package service

import (
	"log/slog"

	"github.com/IslamCHup/coworking-manager-project/internal/models"
	"github.com/IslamCHup/coworking-manager-project/internal/repository"
)

type UserService interface {
	GetUserByID(userID uint) (*models.UserResponseDTO, error)
	UpdateUser(userID uint, req models.UserUpdateDTO) error
	DeleteUser(userID uint) error
	GetAllUsers() ([]models.UserResponseDTO, error)
}

type userService struct {
	repo   repository.UserRepository
	logger *slog.Logger
}

func NewUserService(repo repository.UserRepository,
	logger *slog.Logger,
) UserService {
	return &userService{
		repo:   repo,
		logger: logger,
	}
}

func (s *userService) GetUserByID(userID uint) (*models.UserResponseDTO, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	bookings := make([]models.BookingResDTO, 0, len(user.Bookings))

	for _, b := range user.Bookings {
		bookings = append(bookings, models.BookingResDTO{
			ID:         b.ID,
			UserID:     b.UserID,
			PlaceID:    b.PlaceID,
			StartTime:  b.StartTime,
			EndTime:    b.EndTime,
			TotalPrice: b.TotalPrice,
			Status:     string(b.Status),
		})
	}

	return &models.UserResponseDTO{
		ID:        user.ID,
		Phone:     user.Phone,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Bookings:  bookings,
	}, nil
}

func (s *userService) UpdateUser(userID uint, req models.UserUpdateDTO) error {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		s.logger.Error(
			"user not found",
			"user_id", userID,
			"error", err,
		)
		return err
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}

	if req.LastName != nil {
		user.LastName = *req.LastName
	}

	if err := s.repo.UpdateUser(user); err != nil {
		s.logger.Error(
			"UpdateUser failed",
			"user_id", userID,
			"error", err,
		)
		return err
	}

	s.logger.Info("UpdateUser success", "user_id", userID)
	return nil
}

func (s *userService) DeleteUser(userID uint) error {
	if err := s.repo.DeleteUser(userID); err != nil {
		s.logger.Error(
			"DeleteUser failed",
			"user_id", userID,
			"error", err,
		)
		return err
	}

	s.logger.Info("DeleteUser success", "user_id", userID)
	return nil
}

func (s *userService) GetAllUsers() ([]models.UserResponseDTO, error) {
	users, err := s.repo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	result := make([]models.UserResponseDTO, 0, len(users))

	for _, u := range users {
		result = append(result, models.UserResponseDTO{
			ID:        u.ID,
			Phone:     u.Phone,
			FirstName: u.FirstName,
			LastName:  u.LastName,
		})
	}

	return result, nil
}
