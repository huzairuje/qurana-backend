package models

import (
	"time"
)

type OTP struct {
	ID        uint      `gorm:"primaryKey"`
	Phone     string    `gorm:"size:255;not null;index"` // Index for quick lookup
	Otp       string    `gorm:"size:10;not null"`      // OTPs are usually short
	ExpiredAt time.Time `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
