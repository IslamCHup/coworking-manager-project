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

	User  User  `json:"-"`
	Place Place `json:"-"`
}

type BookingReqDTO struct {
	UserID     uint          `json:"user_id" gorm:"not null;index"`
	PlaceID    uint          `json:"place_id" gorm:"not null;index" binding:"required"`
	StartTime  time.Time     `json:"start_time" gorm:"not null;index" binding:"required"`
	EndTime    time.Time     `json:"end_time" gorm:"not null;index" binding:"required,gtfield=StartTime"`
	TotalPrice float64       `json:"total_price" gorm:"not null"`
	Status     BookingStatus `json:"status" gorm:"not null;default:'active'"`

	User  User  `json:"-"`
	Place Place `json:"-"`
}