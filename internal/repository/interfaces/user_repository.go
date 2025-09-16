package interfaces

import (
	"github.com/AtlasOpx/devprep/internal/models"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uuid.UUID) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	Update(id uuid.UUID, req *models.UpdateProfileRequest) error
	Delete(id uuid.UUID) error
	GetAll() ([]models.User, error)
}
