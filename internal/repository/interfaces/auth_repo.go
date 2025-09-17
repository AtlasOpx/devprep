package interfaces

import (
	"github.com/AtlasOpx/devprep/internal/models"
	"github.com/google/uuid"
	"time"
)

type AuthRepository interface {
	CreateSession(userID uuid.UUID, sessionToken string, expiresAt time.Time, userAgent, ipAddress string) error
	DeleteSession(sessionToken string) error
	GetSessionByToken(sessionToken string) (*models.Session, error)
	ValidateSession(sessionToken string) (*models.User, error)
}
