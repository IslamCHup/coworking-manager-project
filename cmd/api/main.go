package main

import (
	"os"

	"github.com/IslamCHup/coworking-manager-project/internal/config"
	"github.com/IslamCHup/coworking-manager-project/internal/models"
	"github.com/IslamCHup/coworking-manager-project/internal/repository"
	"github.com/IslamCHup/coworking-manager-project/internal/service"
	"github.com/IslamCHup/coworking-manager-project/internal/transport"
	"github.com/gin-gonic/gin"
)

func main() {

	logger := config.InitLogger()

	db := config.SetupDataBase(logger)
	if db == nil {
		logger.Error("Ошибка при установке базы данных: db is nil")
	}

	// Автомиграция моделей
	if err := db.AutoMigrate(
		&models.Admin{},
		&models.User{},
		&models.Booking{},
		&models.Review{},
		&models.Place{},
		&models.PhoneVerification{},
		&models.RefreshToken{},
	); err != nil {
		logger.Error("Ошибка при выполнении автомиграции", "error", err)
		return
	}

	// Инициализация репозиториев
	bookingRepo := repository.NewBookingRepository(db, logger)
	adminRepo := repository.NewAdminRepository(db, logger)
	userRepo := repository.NewUserRepository(db, logger)
	placeRepo := repository.NewPlaceRepository(db, logger)

	// Инициализация сервисов
	bookingService := service.NewBookingService(bookingRepo, placeRepo, db, logger)
	adminService := service.NewAdminService(adminRepo, logger)
	userService := service.NewUserService(userRepo, logger)

	r := gin.Default()

	transport.RegisterRoutes(r, logger, bookingService, adminService, userService)

	logger.Info("Запуск HTTP-сервера", "port", os.Getenv("PORT"))
	if err := r.Run(":" + os.Getenv("PORT")); err != nil {
		logger.Error("не удалось запустить HTTP-сервер", "err", err)
	}
}
