package repository

import (
	"log/slog"

	"github.com/IslamCHup/coworking-manager-project/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByID(id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)

	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	DeleteUser(id uint) error

	GetAllUsers() ([]models.User, error)
	UpdateUserBalance(userID uint, amount int) error
}

type userRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewUserRepository(db *gorm.DB, logger *slog.Logger) UserRepository {
	return &userRepository{
		db:     db,
		logger: logger,
	}
}

func (r *userRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User

	if err := r.db.
		Preload("Bookings").
		Preload("Reviews").
		First(&user, id).Error; err != nil {

		r.logger.Error("GetUserByID failed", "user_id", id, "error", err)
		return nil, err
	}

	r.logger.Info("GetUserByID success", "user_id", user.ID)
	return &user, nil
}

func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	if err := r.db.
		Where("email = ?", email).
		First(&user).Error; err != nil {

		r.logger.Warn("GetUserByEmail failed", "email", email, "error", err)
		return nil, err
	}

	r.logger.Info("GetUserByEmail success", "user_id", user.ID, "email", email)
	return &user, nil
}

func (r *userRepository) CreateUser(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		r.logger.Error(
			"CreateUser failed",
			"email", user.Email,
			"error", err,
		)
		return err
	}

	r.logger.Info(
		"User created",
		"user_id", user.ID,
		"email", user.Email,
	)
	return nil
}

func (r *userRepository) UpdateUser(user *models.User) error {
	if err := r.db.Save(user).Error; err != nil {
		r.logger.Error(
			"UpdateUser failed",
			"user_id", user.ID,
			"error", err,
		)
		return err
	}

	r.logger.Info(
		"User updated",
		"user_id", user.ID,
	)
	return nil
}

func (r *userRepository) DeleteUser(id uint) error {
	res := r.db.Delete(&models.User{}, id)
	if res.Error != nil {
		r.logger.Error("DeleteUser failed", "user_id", id, "error", res.Error)
		return res.Error
	}

	r.logger.Info(
		"DeleteUser success",
		"user_id", id,
		"rows", res.RowsAffected,
	)
	return nil
}

func (r *userRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User

	if err := r.db.Find(&users).Error; err != nil {
		r.logger.Error("GetAllUsers failed", "error", err)
		return nil, err
	}

	r.logger.Info("GetAllUsers success", "count", len(users))
	return users, nil
}

func (r *userRepository) UpdateUserBalance(userID uint, amount int) error {
	result := r.db.
		Model(&models.User{}).
		Where("id = ?", userID).
		Update("balance", gorm.Expr("balance + ?", amount))

	if result.Error != nil {
		r.logger.Error(
			"UpdateUserBalance failed",
			"user_id", userID,
			"amount", amount,
			"error", result.Error,
		)
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.logger.Warn("UpdateUserBalance: user not found", "user_id", userID)
		return gorm.ErrRecordNotFound
	}

	r.logger.Info(
		"UpdateUserBalance success",
		"user_id", userID,
		"amount", amount,
	)
	return nil
}
