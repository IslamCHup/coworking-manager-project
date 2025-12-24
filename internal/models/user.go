package models

type User struct {
	Base

	Phone     string `gorm:"uniqueIndex;not null"`
	FirstName string
	LastName  string
	IsBlocked bool `gorm:"default:false"`
	Balance   int

	Bookings []Booking `gorm:"foreignKey:UserID"`
	Reviews  []Review  `gorm:"foreignKey:UserID"`
}

type UserResponseDTO struct {
	ID        uint   `json:"id"`
	Phone     string `json:"phone"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Balance   int
	Bookings  []BookingResDTO
}

type UserUpdateDTO struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
}
