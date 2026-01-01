package main

import (
	"os"

	"github.com/IslamCHup/coworking-manager-project/internal/config"
	"github.com/IslamCHup/coworking-manager-project/internal/models"

	"github.com/IslamCHup/coworking-manager-project/internal/redis"
	"github.com/IslamCHup/coworking-manager-project/internal/repository"
	"github.com/IslamCHup/coworking-manager-project/internal/service"
	"github.com/IslamCHup/coworking-manager-project/internal/transport"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"net/http"
	_ "net/http/pprof"
)

func main() {
	logger := config.InitLogger()

	if err := godotenv.Load(); err != nil {
		logger.Error("env не найдено")
	}

	db := config.SetupDataBase(logger)
	if db == nil {
		logger.Error("Ошибка при установке базы данных: db is nil")
	}

	redisClient, err := redis.New(os.Getenv("REDIS_ADDR"))
	if err != nil {
		logger.Error("redis unavailable:")
		redisClient = nil
	} else {
		logger.Info("redis connected", "addr", os.Getenv("REDIS_ADDR"))
	}

	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	if err := db.AutoMigrate(
		&models.Admin{},
		&models.User{},
		&models.Booking{},
		&models.Review{},
		&models.Place{},
		&models.RefreshToken{},
	); err != nil {
		logger.Error("Ошибка при выполнении автомиграции", "error", err)
		return
	}

	bookingRepo := repository.NewBookingRepository(db, logger)
	adminRepo := repository.NewAdminRepository(db, logger)
	userRepo := repository.NewUserRepository(db, logger)
	placeRepo := repository.NewPlaceRepository(db, logger)
	refreshRepo := repository.NewRefreshTokenRepository(db, logger)
	reviewRepo := repository.NewReviewRepository(db)

	bookingService := service.NewBookingService(bookingRepo, placeRepo, db, logger, redisClient)
	placeService := service.NewPlaceService(placeRepo, db)
	adminService := service.NewAdminService(adminRepo, logger)
	userService := service.NewUserService(userRepo, logger)
	authService := service.NewAuthService(userRepo, logger)
	refreshService := service.NewRefreshService(refreshRepo, logger)
	reviewService := service.NewReviewService(db, reviewRepo)

	r := gin.Default()

	transport.RegisterRoutes(r, logger, bookingService, placeService, adminService, userService, authService, refreshService, reviewService)

	logger.Info("Запуск HTTP-сервера", "port", os.Getenv("PORT"))
	if err := r.Run(":" + os.Getenv("PORT")); err != nil {
		logger.Error("не удалось запустить HTTP-сервер", "err", err)
	}
}
