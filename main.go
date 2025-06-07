package main

import (
	"log"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"go-htmx-fiber-app/config"
	"go-htmx-fiber-app/database"
	"go-htmx-fiber-app/handlers"
	"go-htmx-fiber-app/middleware" // Import the middleware package
	"github.com/gofiber/template/html/v2" // Import Go template engine
)

func main() {
	appConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	if err := database.Connect(); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// database.Migrate() // Call this manually or via a script if needed, not on every app start for production

	// Initialize Go template engine
	engine := html.New("./views", ".html")
	// Reload templates on each render for development, disable for production
	engine.Reload(true)

	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(logger.New()) // Add logger middleware
	app.Use(cors.New()) // Add CORS middleware for broader compatibility if HTMX requests come from different origins than API

	authHandler := handlers.NewAuthHandler()

	apiV1 := app.Group("/api/v1")
	authRoutes := apiV1.Group("/auth")

	authRoutes.Post("/register", authHandler.Register)
	authRoutes.Post("/login", authHandler.Login)

	// Protected routes
	authRoutes.Post("/verify-otp", middleware.Protected(), authHandler.VerifyOTP)
	authRoutes.Post("/resend-otp", middleware.Protected(), authHandler.ResendOTP)
	authRoutes.Get("/me", middleware.Protected(), authHandler.GetMe)
	authRoutes.Post("/logout", middleware.Protected(), authHandler.Logout) // Logout might also need protection

	// Dashboard/View routes
	viewHandler := handlers.NewViewHandler()
	dashboardRoutes := app.Group("/dashboard")
	dashboardAuthRoutes := dashboardRoutes.Group("/auth")
	dashboardAuthRoutes.Get("/login", viewHandler.ShowLoginPage)
	dashboardAuthRoutes.Get("/register", viewHandler.ShowRegisterPage)
	dashboardAuthRoutes.Get("/verify-otp", viewHandler.ShowOTPVerifyPage) // Page to enter OTP
	dashboardRoutes.Get("/home", viewHandler.ShowDashboardHomePage) // Protected by client-side logic / API calls

	log.Printf("Starting server on port %s", appConfig.ServerPort)
	log.Fatal(app.Listen(":" + appConfig.ServerPort))
}
