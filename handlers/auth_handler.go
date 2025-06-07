package handlers

import (
	"fmt"
	"go-htmx-fiber-app/auth"
	"go-htmx-fiber-app/database" // For GetDB, though repository should abstract this
	"go-htmx-fiber-app/models"
	"go-htmx-fiber-app/repository"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5" // For GetUserIDFromFiberCtx, if used directly here before middleware
)

type AuthHandler struct {
	UserRepo repository.UserRepository
	OTPRepo  repository.OTPRepository
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		UserRepo: repository.UserRepository{},
		OTPRepo:  repository.OTPRepository{},
	}
}

// RegisterRequest defines the expected request body for user registration
type RegisterRequest struct {
	Name     string `json:"name" form:"name"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	Phone    string `json:"phone" form:"phone"` // Assuming phone is a string
}

// LoginRequest defines the expected request body for user login
type LoginRequest struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

// VerifyOTPRequest defines the expected request body for OTP verification
type VerifyOTPRequest struct {
	OTP string `json:"otp" form:"otp"`
}

// Register handles user registration
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	req := new(RegisterRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request"})
	}

	// Basic validation
	if req.Name == "" || req.Email == "" || req.Password == "" || req.Phone == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "All fields (name, email, password, phone) are required"})
	}

	// Check if email or phone already exists
	existingUserByEmail, _ := h.UserRepo.GetUserByEmail(req.Email)
	if existingUserByEmail != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already exists"})
	}
	existingUserByPhone, _ := h.UserRepo.GetUserByPhone(req.Phone)
	if existingUserByPhone != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Phone number already exists"})
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password, // Pass plain password, GORM hook will hash it.
		Phone:    &req.Phone,
	}

	if err := h.UserRepo.CreateUser(user); err != nil {
		log.Printf("Error creating user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to register user"})
	}

	// Generate and store OTP for phone verification
	otp, err := h.OTPRepo.CreateOTP(req.Phone)
	if err != nil {
		log.Printf("Error creating OTP: %v", err)
		// Potentially rollback user creation or handle more gracefully
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create OTP"})
	}
	log.Printf("Generated OTP for %s: %s", req.Phone, otp.Otp) // For testing, normally not logged like this

	// Generate JWT token
	token, err := auth.GenerateJWT(user.ID, user.IsAdmin)
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully. Please verify your phone.",
		"token":   token,
		"user":    user, // Be careful about exposing user data, customize as needed
		"otp_for_testing": otp.Otp, // REMOVE IN PRODUCTION
	})
}

// Login handles user login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	req := new(LoginRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request"})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email and password are required"})
	}

	user, err := h.UserRepo.GetUserByEmail(req.Email)
	if err != nil || user == nil {
		// Important to check user != nil even if err is nil (e.g. gorm.ErrRecordNotFound)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	if !auth.CheckPasswordHash(req.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token, err := auth.GenerateJWT(user.ID, user.IsAdmin)
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to login"})
	}

	phone := ""
    if user.Phone != nil {
        phone = *user.Phone
    }

	return c.JSON(fiber.Map{
		"message": "Logged in successfully",
		"token":   token,
		"user": fiber.Map{
			"id": user.ID,
			"name": user.Name,
			"email": user.Email,
			"phone": phone,
			"is_admin": user.IsAdmin,
			"phone_verified_at": user.PhoneVerifiedAt,
			"email_verified_at": user.EmailVerifiedAt,
		},
	})
}

// VerifyOTP handles OTP verification for the authenticated user
func (h *AuthHandler) VerifyOTP(c *fiber.Ctx) error {
    userID, err := auth.GetUserIDFromFiberCtx(c)
    if err != nil {
        log.Printf("VerifyOTP - Error getting user ID from context: %v. Likely middleware issue or bad token.", err)
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized: Invalid or missing token."})
    }

	req := new(VerifyOTPRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request"})
	}
	if req.OTP == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "OTP is required"})
	}

	user, _ := h.UserRepo.GetUserByID(uint(userID))
	if user == nil || user.Phone == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User or user phone not found"})
	}

	otp, err := h.OTPRepo.GetValidOTPByPhone(*user.Phone, req.OTP)
	if err != nil {
		log.Printf("Error validating OTP: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error validating OTP"})
	}
	if otp == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid or expired OTP"})
	}

	if err := h.UserRepo.UpdateUserPhoneVerificationStatus(user.ID, time.Now()); err != nil {
		log.Printf("Error updating phone verification status: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to verify phone"})
	}

	if err := h.OTPRepo.DeleteOTP(otp); err != nil {
		log.Printf("Error deleting used OTP: %v", err)
	}

	return c.JSON(fiber.Map{"message": "Phone verified successfully"})
}

// ResendOTP handles resending OTP for the authenticated user
func (h *AuthHandler) ResendOTP(c *fiber.Ctx) error {
    userID, err := auth.GetUserIDFromFiberCtx(c)
    if err != nil {
        log.Printf("ResendOTP - Error getting user ID from context: %v", err)
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized: Invalid or missing token."})
    }

	user, _ := h.UserRepo.GetUserByID(uint(userID))
	if user == nil || user.Phone == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User or user phone not found"})
	}

	if user.PhoneVerifiedAt != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Phone already verified"})
	}

	newOtp, err := h.OTPRepo.CreateOTP(*user.Phone)
	if err != nil {
		log.Printf("Error creating new OTP for resend: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to resend OTP"})
	}
	log.Printf("Resent OTP for %s: %s", *user.Phone, newOtp.Otp)

	return c.JSON(fiber.Map{
		"message": "OTP resent successfully",
		"otp_for_testing": newOtp.Otp, // REMOVE IN PRODUCTION
	})
}

// GetMe handles fetching authenticated user's details
func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
    userID, err := auth.GetUserIDFromFiberCtx(c)
    if err != nil {
        log.Printf("GetMe - Error getting user ID from context: %v", err)
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized: Invalid or missing token."})
    }

	user, err := h.UserRepo.GetUserByID(uint(userID))
	if err != nil || user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

    phone := ""
    if user.Phone != nil {
        phone = *user.Phone
    }

	return c.JSON(fiber.Map{
		"id": user.ID,
		"name": user.Name,
		"email": user.Email,
		"phone": phone,
		"is_admin": user.IsAdmin,
		"phone_verified_at": user.PhoneVerifiedAt,
		"email_verified_at": user.EmailVerifiedAt,
	})
}

// Logout handles user logout (currently client-side responsibility for token invalidation)
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Logged out successfully. Please delete your token."})
}
