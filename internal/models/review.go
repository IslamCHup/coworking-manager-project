package models

import "time"

type Review struct {
	Base
	UserID     uint      `json:"user_id" gorm:"not null;index"`
	PlaceID    uint      `json:"place_id" gorm:"not null;index" binding:"required"`
	AdminID    *uint     `json:"admin_id,omitempty" gorm:"index"`
	Rating     int       `json:"rating" gorm:"not null;check:rating >= 1 AND rating <= 5" binding:"required,min=1,max=5"`
	Text       string    `json:"text" gorm:"not null" binding:"required,min=5"`
	IsApproved bool      `json:"is_approved" gorm:"not null;default:false"`
	CreatedAt  time.Time `json:"created_at"`

	User  User   `json:"-"`
	Place Place  `json:"-"`
	Admin *Admin `json:"-"`
}
type PlaceRatingId struct{
	PlaceId uint
	Rating int
	Text string
	CreatedAt time.Time
} 
type PlaceRatingIdUpdate struct{
	PlaceId uint
	Rating int
	Text string
	CreatedAt time.Time
}