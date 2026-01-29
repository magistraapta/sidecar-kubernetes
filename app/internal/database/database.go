package database

import (
	"app/internal/models"
	"log/slog"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		slog.Error("DATABASE_URL environment variable is not set")
		panic("DATABASE_URL environment variable is required")
	}
	
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		slog.Error("failed to connect database", "error", err)
		panic(err)
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		slog.Error("failed to migrate database", "error", err)
		panic(err)
	}

	slog.Info("database connected successfully")
	return db
}
