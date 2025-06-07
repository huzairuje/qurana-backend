package models

import (
	"time"
)

type PasswordResetToken struct {
	Email     string    `gorm:"primaryKey;size:255"`
	Token     string    `gorm:"size:255;uniqueIndex;not null"`
	CreatedAt time.Time
}
