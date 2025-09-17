package routes

import (
	"github.com/AtlasOpx/devprep/internal/app"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes настраивает все маршруты приложения
func SetupRoutes(fiberApp *fiber.App, deps *app.Dependencies) {
	api := fiberApp.Group("/api/v1")

	SetupAuthRoutes(api, deps.AuthHandler, deps.AuthMiddleware)
	SetupUserRoutes(api, deps.UserHandler, deps.AuthMiddleware)
	SetupAdminRoutes(api, deps.UserHandler, deps.AuthMiddleware)
}
