package repository

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"go-htmx-fiber-app/database"
	"go-htmx-fiber-app/models"
	"gorm.io/gorm"
)

type OTPRepository struct{}

// generateRandomOTP generates a 6-digit OTP.
func (r *OTPRepository) generateRandomOTP(length int) (string, error) {
	const digits = "0123456789"
	otp := make([]byte, length)
	for i := range otp {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		otp[i] = digits[num.Int64()]
	}
	return string(otp), nil
}

func (r *OTPRepository) CreateOTP(phone string) (*models.OTP, error) {
	db := database.GetDB()

	otpCode, err := r.generateRandomOTP(6) // 6-digit OTP
	if err != nil {
		return nil, fmt.Errorf("failed to generate OTP code: %w", err)
	}

	otp := &models.OTP{
		Phone:     phone,
		Otp:       otpCode,
		ExpiredAt: time.Now().Add(5 * time.Minute), // OTP expires in 5 minutes
	}
	err = db.Create(otp).Error
	return otp, err
}

// GetValidOTPByPhone retrieves an OTP for a given phone number that has not expired and has not been used.
// In this simplified version, we don't have an explicit 'used' flag, so we rely on deleting/invalidating after use.
func (r *OTPRepository) GetValidOTPByPhone(phone string, otpCode string) (*models.OTP, error) {
	var otp models.OTP
	db := database.GetDB()
	err := db.Where("phone = ? AND otp = ? AND expired_at > ?", phone, otpCode, time.Now()).First(&otp).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil // Not found is a common case, not an error
	}
	return &otp, err
}

func (r *OTPRepository) DeleteOTP(otp *models.OTP) error {
	db := database.GetDB()
	return db.Delete(otp).Error
}
