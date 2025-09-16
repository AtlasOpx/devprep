package service

import (
	"database/sql"
	"github.com/AtlasOpx/devprep/internal/models"
	"github.com/AtlasOpx/devprep/internal/repository/interfaces"
	"github.com/AtlasOpx/devprep/internal/utils"
	"github.com/google/uuid"
	"time"
)

type AuthService struct {
	userRepo interfaces.UserRepository
	authRepo interfaces.AuthRepository
}

func NewAuthService(userRepo interfaces.UserRepository, authRepo interfaces.AuthRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		authRepo: authRepo,
	}
}

func (s *AuthService) Register(req *models.RegisterRequest) (*uuid.UUID, error) {
	existingUser, _ := s.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, sql.ErrNoRows
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	userID := uuid.New()
	user := &models.User{
		ID:           userID,
		Email:        req.Email,
		Username:     req.Username,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		PasswordHash: hashedPassword,
		Role:         models.UserRoleUser,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return &userID, nil
}

func (s *AuthService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, sql.ErrNoRows
	}

	if !user.IsActive {
		return nil, sql.ErrNoRows
	}

	response := &models.LoginResponse{
		Message: "Login successful",
		User:    *user,
	}

	return response, nil
}

func (s *AuthService) Logout(sessionToken string) error {
	return s.authRepo.DeleteSession(sessionToken)
}
