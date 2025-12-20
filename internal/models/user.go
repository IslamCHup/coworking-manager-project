package models

import "time"

type User struct {
	Base

	Email         string    `json:"email" gorm:"uniqueIndex;not null" binding:"required,email"`
	PasswordHash  string    `json:"-" gorm:"not null"`
	FirstName     string    `json:"first_name" gorm:"not null" binding:"required,min=2"`
	LastName      string    `json:"last_name" gorm:"not null" binding:"required,min=2"`
	Balance       float64   `json:"balance" gorm:"not null;default:0"`
	IsBlocked     bool      `json:"-" gorm:"not null;default:false"`
	EmailVerified bool      `json:"-" gorm:"not null;default:false"`
	CreatedAt     time.Time `json:"created_at"`

	Bookings []Booking `json:"-"`
	Reviews  []Review  `json:"-"`
}
