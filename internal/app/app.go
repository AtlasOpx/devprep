package app

import (
	"github.com/AtlasOpx/devprep/internal/config"
	"github.com/AtlasOpx/devprep/internal/database"
	"github.com/AtlasOpx/devprep/internal/handlers"
	"github.com/AtlasOpx/devprep/internal/middleware"
	"github.com/AtlasOpx/devprep/internal/repository"
	"github.com/AtlasOpx/devprep/internal/service"
)

// Dependencies содержит все зависимости приложения
type Dependencies struct {
	AuthHandler    *handlers.AuthHandler
	UserHandler    *handlers.UserHandler
	AuthMiddleware *middleware.AuthMiddleware
}

// NewDependencies создает и инициализирует все зависимости
func NewDependencies(db *database.DB, cfg *config.Config) *Dependencies {
	// Репозитории
	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthRepository(db)

	// Сервисы
	authService := service.NewAuthService(userRepo, authRepo)
	userService := service.NewUserService(userRepo)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService, authRepo, cfg)
	userHandler := handlers.NewUserHandler(userService)

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(authRepo)

	return &Dependencies{
		AuthHandler:    authHandler,
		UserHandler:    userHandler,
		AuthMiddleware: authMiddleware,
	}
}