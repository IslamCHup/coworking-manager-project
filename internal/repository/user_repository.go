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

	if err := r.db.First(&user, id).Error; err != nil {

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
