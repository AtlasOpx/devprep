package dto

import (
	"github.com/AtlasOpx/devprep/internal/models"
)

func RegisterRequestToModel(dto *RegisterRequest) *models.RegisterRequest {
	return &models.RegisterRequest{
		Email:     dto.Email,
		Username:  dto.Username,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Password:  dto.Password,
	}
}

func LoginRequestToModel(dto *LoginRequest) *models.LoginRequest {
	return &models.LoginRequest{
		Email:    dto.Email,
		Password: dto.Password,
	}
}

func UserToDTO(user *models.User) UserDTO {
	return UserDTO{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      string(user.Role),
	}
}

func UserToProfileResponse(user *models.User) UserProfileResponse {
	return UserProfileResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      string(user.Role),
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func UpdateProfileRequestToModel(dto *UpdateProfileRequest) *models.UpdateProfileRequest {
	return &models.UpdateProfileRequest{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Username:  dto.Username,
	}
}

func UserToResponse(user *models.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      string(user.Role),
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func UsersToListResponse(users []models.User) UsersListResponse {
	userDTOs := make([]UserProfileResponse, len(users))
	for i, user := range users {
		userDTOs[i] = UserToProfileResponse(&user)
	}

	return UsersListResponse{
		Users: userDTOs,
		Total: len(users),
	}
}
