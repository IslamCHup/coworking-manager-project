package repository

import (
	"log/slog"

	"github.com/IslamCHup/coworking-manager-project/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByID(id uint) (*models.User, error)
	GetUserByPhone(phone string) (*models.User, error)
	CreateUser(req *models.User) error
	UpdateUser(req *models.User) error
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

	if err := r.db.Preload("Bookings").First(&user, id).Error; err != nil {

		r.logger.Error("GetByID failed", "user_id", id, "error", err)
		return nil, err
	}

	r.logger.Info("GetByID success", "user_id", user.ID)
	return &user, nil

}

func (r *userRepository) GetUserByPhone(phone string) (*models.User, error) {
	var user models.User

	if err := r.db.Where("phone = ?", phone).First(&user).Error; err != nil {
		
			r.logger.Error("GetUserByPhone failed", "phone", phone, "error", err)
		return nil, err
	}

	r.logger.Info("GetUserByPhone success", "user_id", user.ID)
	return &user, nil
}

func (r *userRepository) CreateUser(req *models.User) error {
	if err := r.db.Create(req).Error; err != nil {
		r.logger.Error(
			"failed to create user",
			"phone", req.Phone,
			"error", err,
		)
		return err
	}

	r.logger.Info(
		"user created",
		"user_id", req.ID,
		"phone", req.Phone,
	)
	return nil
}

func (r *userRepository) UpdateUser(req *models.User) error {
	if err := r.db.Save(req).Error; err != nil {
		r.logger.Error(
			"failed to update user",
			"user_id", req.ID,
			"error", err,
		)
		return err
	}

	r.logger.Info(
		"user updated",
		"user_id", req.ID,
	)
	return nil
}

func (r *userRepository) DeleteUser(id uint) error {
	res := r.db.Delete(&models.User{}, id)
	if res.Error != nil {
		r.logger.Error("DeleteUser failed", "id", id, "err", res.Error)
		return res.Error
	}
	r.logger.Info("DeleteUser success", "id", id, "rows", res.RowsAffected)
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
	result := r.db.Model(&models.User{}).Where("id = ?", userID).Update("balance", gorm.Expr("balance + ?", amount))
	if result.Error != nil {
		r.logger.Error("UpdateUserBalance failed", "user_id", userID, "amount", amount, "error", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		r.logger.Warn("UpdateUserBalance: user not found", "user_id", userID)
		return gorm.ErrRecordNotFound
	}
	r.logger.Info("UpdateUserBalance success", "user_id", userID, "amount", amount)
	return nil
}
