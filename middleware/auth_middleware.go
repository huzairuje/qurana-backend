package middleware

import (
	"strings"
	"go-htmx-fiber-app/auth" // Your auth package for JWT validation
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5" // To store the parsed token
)

// Protected is a middleware function to protect routes that require JWT authentication.
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or malformed JWT",
			})
		}

		// Expecting "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Malformed token, expecting Bearer token",
			})
		}
		tokenString := parts[1]

		// Validate the token
		// auth.ValidateJWT now returns (*auth.JWTCustomClaims, error)
		// We need to get the raw token object if GetUserIDFromFiberCtx expects c.Locals("user").(*jwt.Token)
		// Let's adjust ValidateJWT or how we store it.
		// For now, let's parse it again here to get the *jwt.Token or adjust GetUserIDFromFiberCtx.
		// Simpler: auth.ValidateJWT should return the *jwt.Token if successful and that's what we store.
		// Or, GetUserIDFromFiberCtx can take claims directly.

		// Let's assume GetUserIDFromFiberCtx is adapted to take claims, or we store claims.
		// For consistency with its current implementation (expecting *jwt.Token), let's store the token.

		// Re-evaluating: auth.ValidateJWT returns *JWTCustomClaims.
		// auth.GetUserIDFromFiberCtx expects c.Locals("user").(*jwt.Token) then claims := user.Claims.(*JWTCustomClaims)
		// So we need to store the *jwt.Token object itself.
		// The current ValidateJWT in auth/jwt.go parses and returns claims.
		// We'll need to adjust what we store in locals or how GetUserIDFromFiberCtx works.

		// Option 1: Modify ValidateJWT to return the *jwt.Token as well, or just the token if valid.
		// Option 2: Modify GetUserIDFromFiberCtx to accept *JWTCustomClaims from c.Locals.
		// Option 2 is simpler for now.

		claims, err := auth.ValidateJWT(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired JWT",
				"details": err.Error(), // Consider removing detailed errors in prod
			})
		}

		// Store the claims in context for use by route handlers
		c.Locals("user_claims", claims) // Store JWTCustomClaims directly

		return c.Next()
	}
}
