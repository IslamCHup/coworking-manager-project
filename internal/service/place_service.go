package service

import (
	"errors"

	"github.com/IslamCHup/coworking-manager-project/internal/models"
	"github.com/IslamCHup/coworking-manager-project/internal/repository"
	"gorm.io/gorm"
)

var ErrPlaceNotFound = errors.New("place not found")

type PlaceService interface {
	ListPlaces(filter *models.FilterPlace) (*[]models.Place, error)
	GetPlaceByID(id uint) (*models.Place, error)
	ListFreePlaces(filter *models.FilterPlace) (*[]models.Place, error)
}

type placeService struct {
	placeRepo repository.PlaceRepository
	db        *gorm.DB
}

func NewPlaceService(placeRepo repository.PlaceRepository, db *gorm.DB) PlaceService {
	return &placeService{placeRepo: placeRepo, db: db}
}

func (s *placeService) ListPlaces(filter *models.FilterPlace) (*[]models.Place, error) {
	places, err := s.placeRepo.ListPlaces(filter)
	if err != nil {
		return nil, err
	}
	return places, nil
}

func (s *placeService) GetPlaceByID(id uint) (*models.Place, error) {
	place, err := s.placeRepo.GetPlaceByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPlaceNotFound
		}
		return nil, err
	}
	return place, nil
}

func (s *placeService) ListFreePlaces(filter *models.FilterPlace) (*[]models.Place, error) {
	places, err := s.placeRepo.ListFreePlaces(filter)
	if err != nil {
		return nil, err
	}
	return places, nil
}