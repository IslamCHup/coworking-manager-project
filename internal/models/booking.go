package models

import "time"

type BookingStatus string

const (
	BookingNonActive BookingStatus = "non_active"
	BookingActive    BookingStatus = "active"
	BookingCancelled BookingStatus = "cancelled"
)

type Booking struct {
	Base
	UserID     uint          `json:"user_id" gorm:"not null;index"`
	PlaceID    uint          `json:"place_id" gorm:"not null;index" binding:"required"`
	StartTime  time.Time     `json:"start_time" gorm:"not null;index" binding:"required"`
	EndTime    time.Time     `json:"end_time" gorm:"not null;index" binding:"required,gtfield=StartTime"`
	TotalPrice int           `json:"total_price" gorm:"not null"`
	Status     BookingStatus `json:"status" gorm:"not null;default:'non_active'"`

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
	UserID    *uint   `json:"user_id,omitempty"`
	PlaceID   *uint   `json:"place_id,omitempty"`
	StartTime *string `json:"start_time,omitempty"`
	EndTime   *string `json:"end_time,omitempty"`
}

type BookingStatusUpdateDTO struct {
	Status BookingStatus `json:"status" binding:"required,oneof=active cancelled non_active"`
}

type BookingResDTO struct {
	UserID     uint      `json:"user_id"`
	PlaceID    uint      `json:"place_id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	TotalPrice int       `json:"total_price"`
	Status     string    `json:"status"`

	User  *UserResponseDTO `json:"user,omitempty"`
	Place *Place           `json:"place,omitempty"`
}

type FilterBooking struct {
	Status    *string    `form:"status"`
	PriceMin  *int       `form:"price_min"`
	PriceMax  *int       `form:"price_max"`
	StartTime *time.Time `form:"start_time"`
	EndTime   *time.Time `form:"end_time"  `
	Limit     int        `form:"limit"`
	Offset    int        `form:"offset"`
	SortBy    string     `form:"sort_by"`
	Order     string     `form:"order"`
}
