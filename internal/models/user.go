package models

type User struct {
	Base

	Phone     string `gorm:"uniqueIndex;not null"`
	FirstName string
	LastName  string
	IsBlocked bool `gorm:"default:false"`

	Bookings []Booking `json:"-"`
	Reviews  []Review  `json:"-"`
}

type UserResponseDTO struct {
	ID        uint   `json:"id"`
	Phone     string `json:"phone"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UserUpdateDTO struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
}
