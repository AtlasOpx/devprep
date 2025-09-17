package handlers

import (
	"github.com/AtlasOpx/devprep/internal/dto"
	"github.com/AtlasOpx/devprep/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	user, err := h.userService.GetProfile(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{Error: "User not found"})
	}

	response := dto.UserToResponse(user)
	return c.JSON(response)
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	var req dto.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{Error: "Invalid request body"})
	}

	modelReq := dto.UpdateProfileRequestToModel(&req)
	err := h.userService.UpdateProfile(userID, modelReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "Failed to update profile"})
	}

	response := dto.UpdateProfileResponse{Message: "Profile updated successfully"}
	return c.JSON(response)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)

	err := h.userService.DeleteUser(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "Failed to delete user"})
	}

	response := dto.SuccessResponse{Message: "User deleted successfully"}
	return c.JSON(response)
}

func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{Error: "Failed to get users"})
	}

	response := dto.UsersToListResponse(users)
	return c.JSON(response)
}
