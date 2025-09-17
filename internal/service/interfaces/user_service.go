package interfaces

import (
	"github.com/AtlasOpx/devprep/internal/models"
	"github.com/google/uuid"
)

type UserService interface {
	GetProfile(userID uuid.UUID) (*models.User, error)
	UpdateProfile(userID uuid.UUID, req *models.UpdateProfileRequest) error
	GetAllUsers() ([]models.User, error)
}
