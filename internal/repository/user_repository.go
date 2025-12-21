package repository

import "github.com/IslamCHup/coworking-manager-project/internal/models"

type UserRepository interface {
	CreateUser(req *models.User) error
	GetUserByID(id uint) (*models.User, error)
	UpdateUser(req *models.User) error
	DeleteUser(id uint) error
}


