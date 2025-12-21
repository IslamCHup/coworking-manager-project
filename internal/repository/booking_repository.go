package repository

import "github.com/IslamCHup/coworking-manager-project/internal/models"

type BookingRepository interface {
	CreateBooking(req *models.Booking) error
	UpdateBook(req *models.Booking) error
	DeleteBook(id uint) error
	GetBookId(id uint) (*models.Booking, error)
}
