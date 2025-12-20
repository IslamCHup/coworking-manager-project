package main

import (
	"log"
	"os"

	"github.com/IslamCHup/coworking-manager-project/internal/config"
	"github.com/IslamCHup/coworking-manager-project/internal/logging"
	"github.com/gin-gonic/gin"
)

func main() {

	logger := logging.InitLogger()

	db := config.SetupDataBase(logger)

	if err := db.AutoMigrate(); err != nil {
	}

	r := gin.Default()

	if err := r.Run(":" + os.Getenv("PORT")); err != nil {
		log.Fatalf("не удалось запустить HTTP-сервер: %v", err)
	}
}
