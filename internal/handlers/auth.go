package handlers

import (
	"github.com/AtlasOpx/devprep/internal/config"
	"github.com/AtlasOpx/devprep/internal/dto"
	authRepoInterface "github.com/AtlasOpx/devprep/internal/repository/interfaces"
	authServiceInterface "github.com/AtlasOpx/devprep/internal/service/interfaces"
	"github.com/AtlasOpx/devprep/internal/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService authServiceInterface.AuthService
	authRepo    authRepoInterface.AuthRepository
	cfg         *config.Config
}

func NewAuthHandler(authService authServiceInterface.AuthService, authRepo authRepoInterface.AuthRepository, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		authRepo:    authRepo,
		cfg:         cfg,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(dto.ErrorResponse{Error: "Invalid request body"})
	}

	modelReq := dto.RegisterRequestToModel(&req)
	userID, err := h.authService.Register(modelReq)
	if err != nil {
		return c.Status(400).JSON(dto.ErrorResponse{Error: "User already exists or failed to create"})
	}

	response := dto.RegisterResponse{
		Message: "User created successfully",
		UserID:  *userID,
	}

	return c.Status(201).JSON(response)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(dto.ErrorResponse{Error: "Invalid request body"})
	}

	modelReq := dto.LoginRequestToModel(&req)
	response, err := h.authService.Login(modelReq)
	if err != nil {
		return c.Status(401).JSON(dto.ErrorResponse{Error: "Invalid credentials"})
	}

	sessionToken := utils.GenerateSessionToken()
	expiresAt := time.Now().Add(time.Hour * 24)

	err = h.authRepo.CreateSession(response.User.ID, sessionToken, expiresAt, c.Get("User-Agent"), c.IP())
	if err != nil {
		return c.Status(500).JSON(dto.ErrorResponse{Error: "Failed to create session"})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expiresAt,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
	})

	loginResponse := dto.LoginResponse{
		Message: response.Message,
		User:    dto.UserToDTO(&response.User),
	}

	return c.JSON(loginResponse)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	sessionToken := c.Cookies("session_token")
	if sessionToken == "" {
		return c.Status(400).JSON(dto.ErrorResponse{Error: "No session token"})
	}

	err := h.authService.Logout(sessionToken)
	if err != nil {
		return c.Status(500).JSON(dto.ErrorResponse{Error: "Failed to logout"})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	})

	response := dto.LogoutResponse{Message: "Logout successful"}
	return c.JSON(response)
}
