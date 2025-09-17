package service

import (
	"github.com/AtlasOpx/devprep/internal/models"
	"github.com/AtlasOpx/devprep/internal/repository"
	"github.com/google/uuid"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
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

func (s *UserService) GetByID(userID uuid.UUID) (*models.User, error) {
	return s.userRepo.GetByID(userID)
}

func (s *UserService) GetByEmail(email string) (*models.User, error) {
	return s.userRepo.GetByEmail(email)
}

func (s *UserService) GetByUsername(username string) (*models.User, error) {
	return s.userRepo.GetByUsername(username)
}

func (s *UserService) DeleteUser(userID uuid.UUID) error {
	return s.userRepo.Delete(userID)
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.GetAll()
}
