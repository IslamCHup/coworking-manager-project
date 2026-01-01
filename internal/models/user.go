package models

type User struct {
	Base
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`

	Email        string `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string `json:"password_hash" gorm:"not null"`
	IsBlocked    bool   `json:"is_blocked" gorm:"default:false"`
	Balance      int    `json:"balance"`

	Bookings []Booking `json:"bookings" gorm:"foreignKey:UserID"`
	Reviews  []Review  `json:"reviews" gorm:"foreignKey:UserID"`
}
type UserResponseDTO struct {
	ID        uint            `json:"id"`
	Email     string          `json:"email"`
	FirstName string          `json:"first_name"`
	LastName  string          `json:"last_name"`
	Balance   int             `json:"balance"`
	Bookings  []BookingResDTO `json:"bookings,omitempty"`
}

type UserUpdateDTO struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
}

type UpdateBalanceDTO struct {
	Amount int `json:"amount" binding:"required"`
}

type RegisterRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}