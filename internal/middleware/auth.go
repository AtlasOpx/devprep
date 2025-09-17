package middleware

import (
	"fmt"
	"github.com/AtlasOpx/devprep/internal/models"
	authRepoInterface "github.com/AtlasOpx/devprep/internal/repository/interfaces"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	authRepo authRepoInterface.AuthRepository
}

func NewAuthMiddleware(authRepo authRepoInterface.AuthRepository) *AuthMiddleware {
	return &AuthMiddleware{authRepo: authRepo}
}

func (m *AuthMiddleware) RequireAuth(c *fiber.Ctx) error {
	sessionToken := c.Cookies("session_token")
	if sessionToken == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Authentication required"})
	}

	session, err := m.authRepo.GetSessionByToken(sessionToken)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid session"})
	}

	if session.ExpiresAt.Before(time.Now()) {
		err := m.authRepo.DeleteSession(sessionToken)
		if err != nil {
			return fmt.Errorf("couldn't delete the session: %w", err)
		}
		return c.Status(401).JSON(fiber.Map{"error": "Session expired"})
	}

	user, err := m.authRepo.ValidateSession(sessionToken)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid session"})
	}

	c.Locals("user_id", user.ID)
	c.Locals("user_role", user.Role)

	return c.Next()
}

func (m *AuthMiddleware) RequireRole(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("user_role").(models.UserRole)

		if string(userRole) != requiredRole {
			return c.Status(403).JSON(fiber.Map{"error": "Access denied"})
		}

		return c.Next()
	}
}
