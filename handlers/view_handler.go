package handlers

import (
	"github.com/gofiber/fiber/v2"
)

type ViewHandler struct{}

func NewViewHandler() *ViewHandler {
	return &ViewHandler{}
}

func (vh *ViewHandler) ShowLoginPage(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{
		"title": "Login",
	}, "layouts/main")
}

func (vh *ViewHandler) ShowRegisterPage(c *fiber.Ctx) error {
	return c.Render("register", fiber.Map{
		"title": "Register",
	}, "layouts/main")
}

func (vh *ViewHandler) ShowOTPVerifyPage(c *fiber.Ctx) error {
	// Here, you might want to pass the phone number or some indicator to the template
	// For now, it's a generic page.
	return c.Render("otp_verify", fiber.Map{
		"title": "Verify OTP",
	}, "layouts/main")
}

func (vh *ViewHandler) ShowDashboardHomePage(c *fiber.Ctx) error {
    // This page should be protected by a middleware that checks for JWT in localStorage via a client-side redirect,
    // or a server-side middleware that checks a session cookie if you were using them.
    // For a pure HTMX+JWT API, the protection mostly relies on API calls being secured by middleware.Protected().
    // The client-side script in dashboard_home.html attempts a basic check.
	return c.Render("dashboard_home", fiber.Map{
		"title": "Dashboard",
	}, "layouts/main")
}
