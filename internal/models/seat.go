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
