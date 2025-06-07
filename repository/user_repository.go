package repository

import (
	"go-htmx-fiber-app/database"
	"go-htmx-fiber-app/models"
	"gorm.io/gorm"
	"time" // Added missing import
)

type UserRepository struct{}

func (r *UserRepository) CreateUser(user *models.User) error {
	db := database.GetDB()
	return db.Create(user).Error
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	db := database.GetDB()
	err := db.Where("email = ?", email).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil // User not found is not necessarily an application error here
	}
	return &user, err
}

func (r *UserRepository) GetUserByPhone(phone string) (*models.User, error) {
	var user models.User
	db := database.GetDB()
	err := db.Where("phone = ?", phone).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	db := database.GetDB()
	err := db.First(&user, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) UpdateUser(user *models.User) error {
	db := database.GetDB()
	return db.Save(user).Error
}

// Specifically for updating phone verification status
func (r *UserRepository) UpdateUserPhoneVerificationStatus(userID uint, verifiedAt time.Time) error {
    db := database.GetDB()
    return db.Model(&models.User{}).Where("id = ?", userID).Update("phone_verified_at", verifiedAt).Error
}
