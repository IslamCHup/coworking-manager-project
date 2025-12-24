package models

import "time"

type PhoneVerification struct {
	Base

	Phone     string    `gorm:"uniqueIndex;not null"`
	CodeHash  string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
	Attempts  int       `gorm:"default:0"`
}

type PhoneRequestDTO struct {
	Phone string `json:"phone" binding:"required,e164"`
}

type PhoneVerifyDTO struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required"`
}
