package routes

import (
	"github.com/AtlasOpx/devprep/internal/config"
	"github.com/AtlasOpx/devprep/internal/database"
	"github.com/AtlasOpx/devprep/internal/handlers"
	"github.com/AtlasOpx/devprep/internal/middleware"
	"github.com/AtlasOpx/devprep/internal/repository"
	"github.com/AtlasOpx/devprep/internal/service"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, db *database.DB, cfg *config.Config) {
	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthRepository(db)

	authService := service.NewAuthService(userRepo, authRepo)
	userService := service.NewUserService(userRepo)

	authHandler := handlers.NewAuthHandler(authService, authRepo, cfg)
	userHandler := handlers.NewUserHandler(userService)

	authMiddleware := middleware.NewAuthMiddleware(authRepo)

	api := app.Group("/api/v1")

	SetupAuthRoutes(api, authHandler, authMiddleware)
	SetupUserRoutes(api, userHandler, authMiddleware)
	SetupAdminRoutes(api, userHandler, authMiddleware)
}
