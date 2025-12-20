package main

import (
	"log"
	"os"

	"github.com/IslamCHup/coworking-manager-project/internal/config"
	"github.com/gin-gonic/gin"
)

func main() {

	db := config.SetupDataBase()

	if err := db.AutoMigrate(); err != nil {
	}

	r := gin.Default()

	if err := r.Run(":" + os.Getenv("PORT")); err != nil {
		log.Fatalf("не удалось запустить HTTP-сервер: %v", err)
	}
}
