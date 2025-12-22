package main

import (
	"fmt"
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
		logger.Error("database setup failed: db is nil")
	}

	if err := db.AutoMigrate(&models.Booking{}); err != nil {
		logger.Error("auto migrate failed", "error", err)
		panic(fmt.Sprintf("не удалось выполнить миграцию:%v", err))
	}

	bookingRepo := repository.NewBookingRepository(db, logger)

	bookingService := service.NewBookingService(bookingRepo, logger)

	r := gin.Default()

	transport.RegisterRoutes(r, logger, bookingService)

	logger.Info("starting HTTP server", "port", os.Getenv("PORT"))
	if err := r.Run(":" + os.Getenv("PORT")); err != nil {
		logger.Error("не удалось запустить HTTP-сервер", "err", err)
	}
}
