package models

import "time"

type BookingStatus string

const (
	BookingActive    BookingStatus = "active"
	BookingCancelled BookingStatus = "cancelled"
	BookingCompleted BookingStatus = "completed"
)

type Booking struct {
	Base
	UserID     uint          `json:"user_id" gorm:"not null;index"`
	PlaceID    uint          `json:"place_id" gorm:"not null;index" binding:"required"`
	StartTime  time.Time     `json:"start_time" gorm:"not null;index" binding:"required"`
	EndTime    time.Time     `json:"end_time" gorm:"not null;index" binding:"required,gtfield=StartTime"`
	TotalPrice float64       `json:"total_price" gorm:"not null"`
	Status     BookingStatus `json:"status" gorm:"not null;default:'active'"`

	User  *User  `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Place *Place `json:"place,omitempty" gorm:"foreignKey:PlaceID"`
}

type BookingReqDTO struct {
    UserID    uint          `json:"user_id" binding:"required"`
    PlaceID   uint          `json:"place_id" binding:"required"`
    StartTime time.Time     `json:"start_time" binding:"required"`
    EndTime   time.Time     `json:"end_time" binding:"required,gtfield=StartTime"`

    TotalPrice float64       `json:"-" binding:"-"`
    Status     BookingStatus `json:"status"`

    User  *User  `json:"-" binding:"-"`
    Place *Place `json:"-" binding:"-"`
}
