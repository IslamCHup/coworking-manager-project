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
	ListPlaces(filter *models.FilterPlace) (*[]models.Place, error)
	ListFreePlaces(filter *models.FilterPlace) (*[]models.Place, error)
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

func (r *placeRepository) ListPlaces(filter *models.FilterPlace) (*[]models.Place, error) {
	if filter == nil {
		filter = &models.FilterPlace{Limit: 20, Offset: 0, SortBy: "created_at", Order: "asc"}
	}
	if filter.Limit <= 0 || filter.Limit > 100 {
		filter.Limit = 20
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}
	if filter.SortBy == "" {
		filter.SortBy = "created_at"
	}
	if filter.Order == "" {
		filter.Order = "asc"
	}

	var places []models.Place
	query := r.db.Model(&models.Place{})

	if filter.Type != nil {
		query = query.Where("type = ?", *filter.Type)
	}
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	allowed := map[string]bool{"created_at": true, "price_per_hour": true, "id": true}
	sortBy := filter.SortBy
	if !allowed[sortBy] {
		sortBy = "created_at"
	}
	order := filter.Order
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	query = query.Order(sortBy + " " + order).Limit(filter.Limit).Offset(filter.Offset)

	if err := query.Find(&places).Error; err != nil {
		r.logger.Error("ListPlaces failed", "error", err)
		return nil, err
	}
	r.logger.Info("ListPlaces success", "count", len(places))
	return &places, nil
}

// ListFreePlaces возвращает места, на которые нет активной брони в указанном промежутке времени
// Если StartTime и EndTime не заданы, проверяется текущее время
func (r *placeRepository) ListFreePlaces(filter *models.FilterPlace) (*[]models.Place, error) {
	var places []models.Place

	query := r.db.Model(&models.Place{}).Where("is_active = ?", true)

	// постройка подзапроса: существует ли активная бронь, пересекающаяся с заданным периодом
	sub := r.db.Table("bookings").Select("1").Where("bookings.place_id = places.id").Where("status = ?", models.BookingActive)

	if filter != nil && filter.StartTime != nil && filter.EndTime != nil {
		// overlap condition: booking.start_time < end AND booking.end_time > start
		sub = sub.Where("start_time < ? AND end_time > ?", *filter.EndTime, *filter.StartTime)
	} else {
		// проверяем текущее время
		// используем NOW() в SQL
		sub = sub.Where("start_time <= NOW() AND end_time >= NOW()")
	}

	query = query.Where("NOT EXISTS (?)", sub)

	if filter != nil {
		if filter.Limit <= 0 || filter.Limit > 100 {
			filter.Limit = 20
		}
		if filter.Offset < 0 {
			filter.Offset = 0
		}
		query = query.Limit(filter.Limit).Offset(filter.Offset)
	}

	if err := query.Find(&places).Error; err != nil {
		r.logger.Error("ListFreePlaces failed", "error", err)
		return nil, err
	}

	r.logger.Info("ListFreePlaces success", "count", len(places))
	return &places, nil
}
