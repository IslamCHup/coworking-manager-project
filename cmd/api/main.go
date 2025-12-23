package main

import (
	"os"

	"github.com/IslamCHup/coworking-manager-project/internal/config"
	"github.com/IslamCHup/coworking-manager-project/internal/models"
	"github.com/IslamCHup/coworking-manager-project/internal/notification"
	"github.com/IslamCHup/coworking-manager-project/internal/repository"
	"github.com/IslamCHup/coworking-manager-project/internal/service"
	"github.com/IslamCHup/coworking-manager-project/internal/transport"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	logger := config.InitLogger()

	if err := godotenv.Load(); err != nil {
		logger.Error("env не найдено")
	}

	if os.Getenv("JWT_ACCESS_SECRET") == "" {
		logger.Error("JWT_ACCESS_SECRET is not set")
		os.Exit(1)
	}

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

	//подключение к sms aero Ислам, не трогать по братскии
	email := os.Getenv("SMS_AERO_EMAIL")
	apiKey := os.Getenv("SMS_AERO_API_KEY")
	from := os.Getenv("SMS_AERO_FROM")
	if from == "" {
		from = "SMSAero"
	}
	smsSender := notification.NewSmsAeroSender(
		email,
		apiKey,
		from,
		logger,
	)

	// Инициализация репозиториев
	bookingRepo := repository.NewBookingRepository(db, logger)
	adminRepo := repository.NewAdminRepository(db, logger)
	userRepo := repository.NewUserRepository(db, logger)
	placeRepo := repository.NewPlaceRepository(db, logger)
	refreshRepo := repository.NewRefreshTokenRepository(db, logger)
	phoneRepo := repository.NewPhoneVerificationRepository(db, logger)
	reviewRepo := repository.NewReviewRepository(db)

	// Инициализация сервисов
	bookingService := service.NewBookingService(bookingRepo, placeRepo, db, logger)
	adminService := service.NewAdminService(adminRepo, logger)
	userService := service.NewUserService(userRepo, logger)
	authService := service.NewAuthService(phoneRepo, userRepo, logger, smsSender)
	refreshService := service.NewRefreshService(refreshRepo, logger)
	reviewService := service.NewReviewService(db, reviewRepo)

	r := gin.Default()

	transport.RegisterRoutes(r, logger, bookingService, adminService, userService, authService, refreshService, reviewService)

	logger.Info("Запуск HTTP-сервера", "port", os.Getenv("PORT"))
	if err := r.Run(":" + os.Getenv("PORT")); err != nil {
		logger.Error("не удалось запустить HTTP-сервер", "err", err)
	}
}
