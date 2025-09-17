package interfaces

import (
	"github.com/AtlasOpx/devprep/internal/models"
	"github.com/google/uuid"
)

type AuthService interface {
	Register(req *models.RegisterRequest) (*uuid.UUID, error)
	Login(req *models.LoginRequest) (*models.LoginResponse, error)
	Logout(sessionToken string) error
}