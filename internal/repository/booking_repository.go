package repository

import "github.com/IslamCHup/coworking-manager-project/internal/models"

type BookingRepository interface {
	CreateBooking(booking *models.Booking) error
	UpdateBook(book *models.Booking) error
	DeleteBook(id uint) error
	GetBookId(id uint) (*models.Booking, error)
}
