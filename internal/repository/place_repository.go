package repository

import "github.com/IslamCHup/coworking-manager-project/internal/models"

type PlaceRepository interface {
	CreatePlace(req *models.Place) error
	GetPlaceByID(id uint) (*models.Place, error)
	UpdatePlace(req *models.Place) error
	DeletePlace(id uint) error
}