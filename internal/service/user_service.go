package service

import (
	"github.com/AtlasOpx/devprep/internal/models"
	"github.com/AtlasOpx/devprep/internal/repository/interfaces"
	"github.com/google/uuid"
)

type UserService struct {
	userRepo interfaces.UserRepository
}

func NewUserService(userRepo interfaces.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetProfile(userID uuid.UUID) (*models.User, error) {
	return s.userRepo.GetByID(userID)
}

func (s *UserService) UpdateProfile(userID uuid.UUID, req *models.UpdateProfileRequest) error {
	return s.userRepo.Update(userID, req)
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.GetAll()
}
