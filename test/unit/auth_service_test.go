package unit

import (
	"database/sql"
	"testing"
	"time"

	"github.com/AtlasOpx/devprep/internal/models"
	"github.com/AtlasOpx/devprep/internal/service"
	"github.com/AtlasOpx/devprep/internal/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(id uuid.UUID, req *models.UpdateProfileRequest) error {
	args := m.Called(id, req)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) GetAll() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) CreateSession(userID uuid.UUID, token string, expiresAt time.Time, userAgent, ipAddress string) error {
	args := m.Called(userID, token, expiresAt, userAgent, ipAddress)
	return args.Error(0)
}

func (m *MockAuthRepository) GetSessionByToken(token string) (*models.Session, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Session), args.Error(1)
}

func (m *MockAuthRepository) ValidateSession(token string) (*models.User, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthRepository) DeleteSession(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockAuthRepository) CleanupExpiredSessions() error {
	args := m.Called()
	return args.Error(0)
}

func TestAuthService_Register_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockAuthRepo := new(MockAuthRepository)
	authService := service.NewAuthService(mockUserRepo, mockAuthRepo)

	req := &models.RegisterRequest{
		Email:     "test@example.com",
		Username:  "testuser",
		FirstName: "Test",
		LastName:  "User",
		Password:  "password123",
	}

	mockUserRepo.On("GetByEmail", req.Email).Return(nil, sql.ErrNoRows)
	mockUserRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

	userID, err := authService.Register(req)

	assert.NoError(t, err)
	assert.NotNil(t, userID)
	mockUserRepo.AssertExpectations(t)
}

func TestAuthService_Register_UserAlreadyExists(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockAuthRepo := new(MockAuthRepository)
	authService := service.NewAuthService(mockUserRepo, mockAuthRepo)

	req := &models.RegisterRequest{
		Email:     "test@example.com",
		Username:  "testuser",
		FirstName: "Test",
		LastName:  "User",
		Password:  "password123",
	}

	existingUser := &models.User{
		ID:    uuid.New(),
		Email: req.Email,
	}

	mockUserRepo.On("GetByEmail", req.Email).Return(existingUser, nil)

	userID, err := authService.Register(req)

	assert.Error(t, err)
	assert.Nil(t, userID)
	assert.Equal(t, sql.ErrNoRows, err)
	mockUserRepo.AssertExpectations(t)
}

func TestAuthService_Login_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockAuthRepo := new(MockAuthRepository)
	authService := service.NewAuthService(mockUserRepo, mockAuthRepo)

	req := &models.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	hashedPassword, _ := utils.HashPassword("password123")

	user := &models.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: hashedPassword,
		IsActive:     true,
	}

	mockUserRepo.On("GetByEmail", req.Email).Return(user, nil)

	response, err := authService.Login(req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "Login successful", response.Message)
	assert.Equal(t, user.ID, response.User.ID)
	mockUserRepo.AssertExpectations(t)
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockAuthRepo := new(MockAuthRepository)
	authService := service.NewAuthService(mockUserRepo, mockAuthRepo)

	req := &models.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	hashedPassword, _ := utils.HashPassword("password123")

	user := &models.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: hashedPassword,
		IsActive:     true,
	}

	mockUserRepo.On("GetByEmail", req.Email).Return(user, nil)

	response, err := authService.Login(req)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, sql.ErrNoRows, err)
	mockUserRepo.AssertExpectations(t)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockAuthRepo := new(MockAuthRepository)
	authService := service.NewAuthService(mockUserRepo, mockAuthRepo)

	req := &models.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	mockUserRepo.On("GetByEmail", req.Email).Return(nil, sql.ErrNoRows)

	response, err := authService.Login(req)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, sql.ErrNoRows, err)
	mockUserRepo.AssertExpectations(t)
}

func TestAuthService_Login_InactiveUser(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockAuthRepo := new(MockAuthRepository)
	authService := service.NewAuthService(mockUserRepo, mockAuthRepo)

	req := &models.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	hashedPassword, _ := utils.HashPassword("password123")

	user := &models.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: hashedPassword,
		IsActive:     false,
	}

	mockUserRepo.On("GetByEmail", req.Email).Return(user, nil)

	response, err := authService.Login(req)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, sql.ErrNoRows, err)
	mockUserRepo.AssertExpectations(t)
}

func TestAuthService_Logout_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockAuthRepo := new(MockAuthRepository)
	authService := service.NewAuthService(mockUserRepo, mockAuthRepo)

	sessionToken := "valid-session-token"

	mockAuthRepo.On("DeleteSession", sessionToken).Return(nil)

	err := authService.Logout(sessionToken)

	assert.NoError(t, err)
	mockAuthRepo.AssertExpectations(t)
}
