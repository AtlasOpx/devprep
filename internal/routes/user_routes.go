package routes

import (
	"github.com/AtlasOpx/devprep/internal/handlers"
	"github.com/AtlasOpx/devprep/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(api fiber.Router, userHandler *handlers.UserHandler, authMiddleware *middleware.AuthMiddleware) {
	user := api.Group("/user")
	user.Use(authMiddleware.RequireAuth)

	user.Get("/profile", userHandler.GetProfile)
	user.Put("/profile", userHandler.UpdateProfile)
}
