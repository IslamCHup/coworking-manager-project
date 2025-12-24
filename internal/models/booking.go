package models

import "time"

type BookingStatus string

const (
	BookingNonActive BookingStatus = "non_active"
	BookingActive    BookingStatus = "active"
	BookingCancelled BookingStatus = "cancelled"
)

const PriceHour = 100

type Booking struct {
	Base
	UserID     uint          `json:"user_id" gorm:"not null;index"`
	PlaceID    uint          `json:"place_id" gorm:"not null;index" binding:"required"`
	StartTime  time.Time     `json:"start_time" gorm:"not null;index" binding:"required"`
	EndTime    time.Time     `json:"end_time" gorm:"not null;index" binding:"required,gtfield=StartTime"`
	TotalPrice float64       `json:"total_price" gorm:"not null"`
	Status     BookingStatus `json:"status" gorm:"not null;default:'non_active'" binding:"oneof=active cancelled non_active"`

	User  *User  `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Place *Place `json:"place,omitempty" gorm:"foreignKey:PlaceID"`
}

type BookingReqDTO struct {
	UserID    uint   `json:"user_id"`
	PlaceID   uint   `json:"place_id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type BookingReqUpdateDTO struct {
	UserID    *uint          `json:"user_id,omitempty"`
	PlaceID   *uint          `json:"place_id,omitempty"`
	StartTime *string        `json:"start_time,omitempty"`
	EndTime   *string        `json:"end_time,omitempty"`
}

type BookingStatusDTO struct{
	Status    BookingStatus `json:"status,omitempty" binding:"oneof=active cancelled non_active"` 
}

type BookingResDTO struct {
	UserID     uint      `json:"user_id"`
	PlaceID    uint      `json:"place_id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"`

	User  *UserResponseDTO `json:"user,omitempty"`
	Place *Place           `json:"place,omitempty"`
}

type BookingUpdateDTO struct {
	UserID    *uint          `json:"user_id"`
	PlaceID   *uint          `json:"place_id"`
	StartTime *time.Time     `json:"start_time"`
	EndTime   *time.Time     `json:"end_time"`
	Status    *BookingStatus `json:"status"`
}

type FilterBooking struct {
	Status    *string    `form:"status" binding:"oneof=active cancelled non_active"`
	PriceMin  *float64   `form:"price_min"`
	PriceMax  *float64   `form:"price_max"`
	StartTime *time.Time `form:"start_time"`
	EndTime   *time.Time `form:"end_time"  `
	Limit     int        `form:"limit"`
	Offset    int        `form:"offset"`
	SortBy    string     `form:"sort_by"`
	Order     string     `form:"order"`
}
