package repository

import "github.com/IslamCHup/coworking-manager-project/internal/models"

type AdminRepository interface {
	CreateAdmin(req *models.Admin) error
	GetAdminByID(id uint) (*models.Admin, error)
	UpdateAdmin(req *models.Admin) error
	DeleteAdmin(id uint) error
}
