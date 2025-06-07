package database

import (
	"fmt"
	"log"

	"go-htmx-fiber-app/config" // Add this import
	"go-htmx-fiber-app/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Migrate() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration for migration: %v", err)
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database for migration: %v", err)
	}

	log.Println("Running database migrations...")

	err = db.AutoMigrate(
		&models.User{},
		&models.OTP{},
		&models.PasswordResetToken{},
	)

	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database migration completed successfully!")
}
