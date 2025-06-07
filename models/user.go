package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID                uint           `gorm:"primaryKey"`
	Name              string         `gorm:"size:255;not null"`
	IsAdmin           bool           `gorm:"default:false"`
	Email             string         `gorm:"size:255;uniqueIndex;not null"`
	EmailVerifiedAt   *time.Time
	Phone             *string        `gorm:"size:255;uniqueIndex"` // Assuming phone should also be unique if present
	PhoneVerifiedAt   *time.Time
	Password          string         `gorm:"size:255;not null"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// BeforeCreate GORM hook to hash the password
func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return
}

// Helper function to check password (not part of GORM model but useful)
func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}
