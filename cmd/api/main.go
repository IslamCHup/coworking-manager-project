package main

import (
	"os"

	"github.com/IslamCHup/coworking-manager-project/internal/config"
	"github.com/gin-gonic/gin"
)

func main() {

	logger := config.InitLogger()

	db := config.SetupDataBase(logger)
	if db == nil {
		logger.Error("database setup failed: db is nil")
	}

	if err := db.AutoMigrate(); err != nil {
		logger.Error("auto migrate failed", "error", err)
	}

	r := gin.Default()

	logger.Info("starting HTTP server", "port", os.Getenv("PORT"))
	if err := r.Run(":" + os.Getenv("PORT")); err != nil {
		logger.Error("не удалось запустить HTTP-сервер", "err", err)
	}
}
