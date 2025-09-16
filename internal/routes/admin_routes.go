package routes

import (
	"github.com/AtlasOpx/devprep/internal/handlers"
	"github.com/AtlasOpx/devprep/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupAdminRoutes(api fiber.Router, userHandler *handlers.UserHandler, authMiddleware *middleware.AuthMiddleware) {
	admin := api.Group("/admin")
	admin.Use(authMiddleware.RequireAuth)
	admin.Use(authMiddleware.RequireRole("admin"))

	admin.Get("/users", userHandler.GetAllUsers)
	//admin.Get("/users/:id", userHandler.GetUserByID)
	//admin.Put("/users/:id", userHandler.UpdateUserByID)
	//admin.Delete("/users/:id", userHandler.DeleteUserByID)
	//admin.Post("/users/:id/ban", userHandler.BanUser)
	//admin.Post("/users/:id/unban", userHandler.UnbanUser)
	//
	//admin.Get("/stats", userHandler.GetSystemStats)
}
