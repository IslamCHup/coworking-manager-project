package models

import "time"

type PlaceType string

const (
	PlaceWorkspace   PlaceType = "workspace"
	PlaceMeetingRoom PlaceType = "meeting_room"
)

type Place struct {
	Base

	Name         string    `json:"name" gorm:"not null" binding:"required,min=2"`
	Type         PlaceType `json:"type" gorm:"not null" binding:"required,oneof=workspace meeting_room"`
	Description  string    `json:"description"`
	PricePerHour int       `json:"price_per_hour" gorm:"not null" binding:"required,gt=0"` // в копейках
	IsActive     bool      `json:"is_active" gorm:"not null;default:true"`
	CreatedAt    time.Time `json:"created_at"`

	Bookings []Booking `json:"-"`
	Reviews  []Review  `json:"-"`
}

// FilterPlace используется для листинга мест и поиска свободных мест
type FilterPlace struct {
	Type      *string    `form:"type" binding:"omitempty,oneof=workspace meeting_room"`
	IsActive  *bool      `form:"is_active"`
	StartTime *time.Time `form:"start_time"`
	EndTime   *time.Time `form:"end_time"`
	Limit     int        `form:"limit"`
	Offset    int        `form:"offset"`
	SortBy    string     `form:"sort_by"`
	Order     string     `form:"order"`
}
