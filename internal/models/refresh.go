package models

import "time"

type RefreshToken struct {
	Base

	UserID    uint      `gorm:"not null;uniqueIndex"`
	TokenHash string    `gorm:"not null;uniqueIndex"`
	ExpiresAt time.Time `gorm:"not null"`
}

type RefreshRequestDTO struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AuthResponseDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
