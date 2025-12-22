package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDataBase(logger *slog.Logger) *gorm.DB {

	if err := godotenv.Load(".env"); err != nil {
		logger.Warn("failed to load .env, proceeding with environment variables", "err", err)
	} else {
		logger.Debug("loaded .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbSSL := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v", dbHost, dbUser, dbPass, dbName, dbPort, dbSSL)

	logger.Debug("prepared database DSN")
	logger.Info("server starting", slog.String("addr", ":"+dbPort), slog.String("env", "local"))
	logger.Debug("opening database connection", slog.String("host", dbHost), slog.String("db", dbName))

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, 
	}), &gorm.Config{})

	if err != nil {
		logger.Error("failed to connect to database", "err", err)
		panic(err)
	}

	logger.Info("connected to database", slog.String("host", dbHost), slog.String("name", dbName), slog.String("port", dbPort))
	return db
}
