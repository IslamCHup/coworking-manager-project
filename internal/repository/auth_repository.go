package repository

import (
	"log/slog"

	"github.com/IslamCHup/coworking-manager-project/internal/models"
	"gorm.io/gorm"
)

type AuthRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
}

type authRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewAuthRepository(
	db *gorm.DB,
	logger *slog.Logger,
) AuthRepository {
	return &authRepository{
		db:     db,
		logger: logger,
	}
}

func (r *authRepository) CreateUser(user *models.User) error {
	err := r.db.Create(user).Error
	if err != nil {
		r.logger.Error(
			"CreateUser failed",
			"email", user.Email,
			"error", err,
		)
		return err
	}

	r.logger.Info(
		"User created",
		"email", user.Email,
	)
	return nil
}

func (r *authRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	err := r.db.
		Where("email = ?", email).
		First(&user).Error

	if err != nil {
		r.logger.Warn(
			"User not found",
			"email", email,
			"error", err,
		)
		return nil, err
	}

	return &user, nil
}
