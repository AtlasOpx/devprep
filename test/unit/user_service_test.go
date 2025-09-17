package unit

import (
	"database/sql"
	"testing"

	"github.com/AtlasOpx/devprep/internal/models"
	"github.com/AtlasOpx/devprep/internal/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserService_GetByID_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	userService := service.NewUserService(mockUserRepo)

	userID := uuid.New()
	expectedUser := &models.User{
		ID:        userID,
		Email:     "test@example.com",
		Username:  "testuser",
		FirstName: "Test",
		LastName:  "User",
		IsActive:  true,
	}

	mockUserRepo.On("GetByID", userID).Return(expectedUser, nil)

	user, err := userService.GetByID(userID)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Email, user.Email)
	mockUserRepo.AssertExpectations(t)
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	userService := service.NewUserService(mockUserRepo)

	userID := uuid.New()

	mockUserRepo.On("GetByID", userID).Return(nil, sql.ErrNoRows)

	user, err := userService.GetByID(userID)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, sql.ErrNoRows, err)
	mockUserRepo.AssertExpectations(t)
}

func TestUserService_UpdateProfile_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	userService := service.NewUserService(mockUserRepo)

	userID := uuid.New()
	req := &models.UpdateProfileRequest{
		FirstName: "Updated",
		LastName:  "Name",
		Username:  "newusername",
	}

	mockUserRepo.On("Update", userID, req).Return(nil)

	err := userService.UpdateProfile(userID, req)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

func TestUserService_DeleteUser_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	userService := service.NewUserService(mockUserRepo)

	userID := uuid.New()

	mockUserRepo.On("Delete", userID).Return(nil)

	err := userService.DeleteUser(userID)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

func TestUserService_GetAllUsers_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	userService := service.NewUserService(mockUserRepo)

	expectedUsers := []models.User{
		{
			ID:        uuid.New(),
			Email:     "user1@example.com",
			Username:  "user1",
			FirstName: "User",
			LastName:  "One",
			IsActive:  true,
		},
		{
			ID:        uuid.New(),
			Email:     "user2@example.com",
			Username:  "user2",
			FirstName: "User",
			LastName:  "Two",
			IsActive:  true,
		},
	}

	mockUserRepo.On("GetAll").Return(expectedUsers, nil)

	users, err := userService.GetAllUsers()

	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, users, 2)
	assert.Equal(t, expectedUsers[0].Email, users[0].Email)
	assert.Equal(t, expectedUsers[1].Email, users[1].Email)
	mockUserRepo.AssertExpectations(t)
}

func TestUserService_GetByEmail_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	userService := service.NewUserService(mockUserRepo)

	email := "test@example.com"
	expectedUser := &models.User{
		ID:        uuid.New(),
		Email:     email,
		Username:  "testuser",
		FirstName: "Test",
		LastName:  "User",
		IsActive:  true,
	}

	mockUserRepo.On("GetByEmail", email).Return(expectedUser, nil)

	user, err := userService.GetByEmail(email)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.Email, user.Email)
	mockUserRepo.AssertExpectations(t)
}

func TestUserService_GetByUsername_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	userService := service.NewUserService(mockUserRepo)

	username := "testuser"
	expectedUser := &models.User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		Username:  username,
		FirstName: "Test",
		LastName:  "User",
		IsActive:  true,
	}

	mockUserRepo.On("GetByUsername", username).Return(expectedUser, nil)

	user, err := userService.GetByUsername(username)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.Username, user.Username)
	mockUserRepo.AssertExpectations(t)
}
