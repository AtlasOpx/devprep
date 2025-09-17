// internal/routes/auth_routes.go
package routes

import (
	"github.com/AtlasOpx/devprep/internal/handlers"
	"github.com/AtlasOpx/devprep/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(api fiber.Router, authHandler *handlers.AuthHandler, authMiddleware *middleware.AuthMiddleware) {
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/logout", authMiddleware.RequireAuth, authHandler.Logout)
}
