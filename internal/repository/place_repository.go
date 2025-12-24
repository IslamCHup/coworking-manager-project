package repository

import (
	"log/slog"

	"github.com/IslamCHup/coworking-manager-project/internal/models"
	"gorm.io/gorm"
)

type PlaceRepository interface {
	CreatePlace(req *models.Place) error
	GetPlaceByID(id uint) (*models.Place, error)
	UpdatePlace(req *models.Place) error
	DeletePlace(id uint) error
}

type placeRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewPlaceRepository(db *gorm.DB, logger *slog.Logger) PlaceRepository {
	return &placeRepository{
		db:     db,
		logger: logger,
	}
}

func (r *placeRepository) CreatePlace(req *models.Place) error {
	if err := r.db.Create(req).Error; err != nil {
		r.logger.Error("CreatePlace failed", "error", err)
		return err
	}
	r.logger.Info("Place created", "place_id", req.ID)
	return nil
}

func (r *placeRepository) GetPlaceByID(id uint) (*models.Place, error) {
	var place models.Place
	if err := r.db.First(&place, id).Error; err != nil {
		r.logger.Error("GetPlaceByID failed", "place_id", id, "error", err)
		return nil, err
	}
	r.logger.Info("GetPlaceByID success", "place_id", id)
	return &place, nil
}

func (r *placeRepository) UpdatePlace(req *models.Place) error {
	if err := r.db.Save(req).Error; err != nil {
		r.logger.Error("UpdatePlace failed", "place_id", req.ID, "error", err)
		return err
	}
	r.logger.Info("Place updated", "place_id", req.ID)
	return nil
}

func (r *placeRepository) DeletePlace(id uint) error {
	if err := r.db.Delete(&models.Place{}, id).Error; err != nil {
		r.logger.Error("DeletePlace failed", "place_id", id, "error", err)
		return err
	}
	r.logger.Info("Place deleted", "place_id", id)
	return nil
}
