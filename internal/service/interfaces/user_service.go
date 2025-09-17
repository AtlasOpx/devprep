package interfaces

import (
	"github.com/AtlasOpx/devprep/internal/models"
	"github.com/google/uuid"
)

type UserService interface {
	GetProfile(userID uuid.UUID) (*models.User, error)
	GetByID(userID uuid.UUID) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	UpdateProfile(userID uuid.UUID, req *models.UpdateProfileRequest) error
	DeleteUser(userID uuid.UUID) error
	GetAllUsers() ([]models.User, error)
}
