package repository

import (
	"log/slog"

	"github.com/IslamCHup/coworking-manager-project/internal/models"
	"gorm.io/gorm"
)

type AdminRepository interface {
	GetAdminByLogin(login string) (*models.Admin, error)
	CreateAdmin(admin *models.Admin) error
	UpdateAdmin(admin *models.Admin) error
}

type adminRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewAdminRepository(db *gorm.DB, logger *slog.Logger) AdminRepository {
	return &adminRepository{
		db:     db,
		logger: logger,
	}
}

func (r *adminRepository) GetAdminByLogin(login string) (*models.Admin, error) {
	var admin models.Admin
	if err := r.db.Where("login = ?", login).First(&admin).Error; err != nil {
		r.logger.Warn("GetAdminByLogin failed", "login", login, "error", err)
		return nil, err
	}
	r.logger.Info("GetAdminByLogin success", "admin_id", admin.ID, "login", login)
	return &admin, nil
}

func (r *adminRepository) CreateAdmin(admin *models.Admin) error {
	if err := r.db.Create(admin).Error; err != nil {
		r.logger.Error("CreateAdmin failed", "error", err)
		return err
	}
	r.logger.Info("CreateAdmin success", "admin_id", admin.ID, "login", admin.Login)
	return nil
}

func (r *adminRepository) UpdateAdmin(admin *models.Admin) error {
	if err := r.db.Save(admin).Error; err != nil {
		r.logger.Error("UpdateAdmin failed", "error", err)
		return err
	}
	r.logger.Info("UpdateAdmin success", "admin_id", admin.ID, "login", admin.Login)
	return nil
}
