package database

import (
	"fmt"
	"log"
	"sync"

	"go-htmx-fiber-app/config"
	"go-htmx-fiber-app/models" // Required for AutoMigrate if not done in a separate migrate package
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB   *gorm.DB
	once sync.Once
)

// Connect initializes the database connection
func Connect() error {
	var err error
	once.Do(func() {
		cfg := config.AppConfig // Use loaded AppConfig
		if cfg == nil {
			log.Fatal("Configuration not loaded. Call config.LoadConfig() first.")
			// Or handle error more gracefully
			err = fmt.Errorf("configuration not loaded")
			return
		}

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
			cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
			return // err will be captured by the outer err
		}

		log.Println("Database connection established successfully.")

		// Optional: AutoMigrate here if you want Connect to also handle migrations
		// Or keep migrations separate as currently in database/migrate.go
		// For this example, let's assume migrations are handled separately or already run.
		// If you want to migrate here:
		// err = DB.AutoMigrate(&models.User{}, &models.OTP{}, &models.PasswordResetToken{})
		// if err != nil {
		//    log.Fatalf("Failed to auto-migrate database: %v", err)
		// }
	})
	return err
}

// GetDB returns the singleton database instance
func GetDB() *gorm.DB {
	if DB == nil {
		// This might happen if Connect() failed or was not called.
		// Consider how to handle this scenario based on application needs.
		// For now, log a fatal error.
		log.Fatal("Database instance is not initialized. Call database.Connect() first.")
	}
	return DB
}
