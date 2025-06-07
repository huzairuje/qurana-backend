package auth

import (
	"fmt"
	"time"
	"strconv"

	"go-htmx-fiber-app/config"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// JWTCustomClaims defines the custom claims for the JWT
type JWTCustomClaims struct {
	UserID  string `json:"user_id"`
	IsAdmin bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

// GenerateJWT generates a new JWT token for a given user ID and admin status.
func GenerateJWT(userID uint, isAdmin bool) (string, error) {
	cfg := config.AppConfig
	if cfg == nil {
		return "", fmt.Errorf("configuration not loaded")
	}

	claims := JWTCustomClaims{
		UserID:  strconv.FormatUint(uint64(userID), 10),
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)), // Token expires in 72 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "go-htmx-fiber-app",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return signedToken, nil
}

// ValidateJWT validates the given token string.
// Returns the claims if the token is valid, otherwise returns an error.
func ValidateJWT(tokenString string) (*JWTCustomClaims, error) {
	cfg := config.AppConfig
	if cfg == nil {
		return nil, fmt.Errorf("configuration not loaded")
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*JWTCustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

// GetUserIDFromFiberCtx extracts UserID from JWT token in Fiber context
func GetUserIDFromFiberCtx(c *fiber.Ctx) (uint, error) {
    claims, ok := c.Locals("user_claims").(*JWTCustomClaims)
    if !ok || claims == nil {
        return 0, fmt.Errorf("user claims not found in context or type assertion failed")
    }

    userID, err := strconv.ParseUint(claims.UserID, 10, 64) // Use 64 for ParseUint
    if err != nil {
        return 0, fmt.Errorf("failed to parse userID from token claims: %w", err)
    }
    return uint(userID), nil
}
