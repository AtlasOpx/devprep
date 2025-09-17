package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AtlasOpx/devprep/internal/config"
	"github.com/AtlasOpx/devprep/internal/database"
	"github.com/AtlasOpx/devprep/internal/dto"
	"github.com/AtlasOpx/devprep/internal/handlers"
	"github.com/AtlasOpx/devprep/internal/middleware"
	"github.com/AtlasOpx/devprep/internal/repository"
	"github.com/AtlasOpx/devprep/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserHandlerTestSuite struct {
	suite.Suite
	app         *fiber.App
	db          *database.DB
	cfg         *config.Config
	userRepo    *repository.UserRepository
	authRepo    *repository.AuthRepository
	userService *service.UserService
	authService *service.AuthService
	userHandler *handlers.UserHandler
	authHandler *handlers.AuthHandler
}

func (suite *UserHandlerTestSuite) SetupSuite() {
	suite.cfg = &config.Config{
		DatabaseURL: "postgres://postgres:password@localhost:5432/devprep_test?sslmode=disable",
		ServerPort:  "3001",
	}

	var err error
	suite.db, err = database.Connect(suite.cfg)
	if err != nil {
		suite.T().Skipf("Database connection failed: %v", err)
	}

	suite.userRepo = repository.NewUserRepository(suite.db).(*repository.UserRepository)
	suite.authRepo = repository.NewAuthRepository(suite.db).(*repository.AuthRepository)
	suite.userService = service.NewUserService(suite.userRepo)
	suite.authService = service.NewAuthService(suite.userRepo, suite.authRepo)
	suite.userHandler = handlers.NewUserHandler(suite.userService)
	suite.authHandler = handlers.NewAuthHandler(suite.authService, suite.authRepo, suite.cfg)

	suite.app = fiber.New()
	suite.setupRoutes()
}

func (suite *UserHandlerTestSuite) TearDownSuite() {
	if suite.db != nil {
		suite.db.Close()
	}
}

func (suite *UserHandlerTestSuite) SetupTest() {
	suite.cleanupDatabase()
}

func (suite *UserHandlerTestSuite) TearDownTest() {
	suite.cleanupDatabase()
}

func (suite *UserHandlerTestSuite) cleanupDatabase() {
	if suite.db != nil {
		_, _ = suite.db.DB.Exec("DELETE FROM sessions")
		_, _ = suite.db.DB.Exec("DELETE FROM users")
	}
}

func (suite *UserHandlerTestSuite) setupRoutes() {
	api := suite.app.Group("/api/v1")

	auth := api.Group("/auth")
	auth.Post("/register", suite.authHandler.Register)
	auth.Post("/login", suite.authHandler.Login)

	authMiddleware := middleware.NewAuthMiddleware(suite.authRepo)

	users := api.Group("/users")
	users.Use(authMiddleware.RequireAuth)
	users.Get("/profile", suite.userHandler.GetProfile)
	users.Put("/profile", suite.userHandler.UpdateProfile)
	users.Delete("/profile", suite.userHandler.DeleteUser)
	users.Get("/", suite.userHandler.GetAllUsers)
}

func (suite *UserHandlerTestSuite) createUserAndLogin() (string, uuid.UUID) {
	registerReq := dto.RegisterRequest{
		Email:     "test@example.com",
		Username:  "testuser",
		FirstName: "Test",
		LastName:  "User",
		Password:  "password123",
	}

	regBody, _ := json.Marshal(registerReq)
	regReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(regBody))
	regReq.Header.Set("Content-Type", "application/json")
	regResp, _ := suite.app.Test(regReq)

	var regResponse dto.RegisterResponse
	err := json.NewDecoder(regResp.Body).Decode(&regResponse)
	if err != nil {
		return "", [16]byte{}
	}

	loginReq := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	loginBody, _ := json.Marshal(loginReq)
	loginReqHttp := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(loginBody))
	loginReqHttp.Header.Set("Content-Type", "application/json")
	loginResp, _ := suite.app.Test(loginReqHttp)

	sessionToken := extractSessionTokenFromCookie(loginResp.Header.Get("Set-Cookie"))

	return sessionToken, regResponse.UserID
}

func (suite *UserHandlerTestSuite) TestGetProfile_Success() {
	sessionToken, userID := suite.createUserAndLogin()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/profile", nil)
	req.Header.Set("Cookie", "session_token="+sessionToken)

	resp, err := suite.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response dto.UserResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return
	}
	assert.Equal(suite.T(), userID, response.ID)
	assert.Equal(suite.T(), "test@example.com", response.Email)
	assert.Equal(suite.T(), "testuser", response.Username)
}

func (suite *UserHandlerTestSuite) TestGetProfile_Unauthorized() {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/profile", nil)

	resp, err := suite.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
}

func (suite *UserHandlerTestSuite) TestUpdateProfile_Success() {
	sessionToken, _ := suite.createUserAndLogin()

	updateReq := dto.UpdateProfileRequest{
		FirstName: "Updated",
		LastName:  "Name",
		Username:  "updateduser",
	}

	reqBody, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/profile", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "session_token="+sessionToken)

	resp, err := suite.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response dto.UpdateProfileResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return
	}
	assert.Equal(suite.T(), "Profile updated successfully", response.Message)
}

func (suite *UserHandlerTestSuite) TestUpdateProfile_Unauthorized() {
	updateReq := dto.UpdateProfileRequest{
		FirstName: "Updated",
		LastName:  "Name",
	}

	reqBody, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/profile", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
}

func (suite *UserHandlerTestSuite) TestDeleteUser_Success() {
	sessionToken, _ := suite.createUserAndLogin()

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/profile", nil)
	req.Header.Set("Cookie", "session_token="+sessionToken)

	resp, err := suite.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response dto.SuccessResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return
	}
	assert.Equal(suite.T(), "User deleted successfully", response.Message)
}

func (suite *UserHandlerTestSuite) TestGetAllUsers_Success() {
	sessionToken, _ := suite.createUserAndLogin()

	registerReq2 := dto.RegisterRequest{
		Email:     "test2@example.com",
		Username:  "testuser2",
		FirstName: "Test2",
		LastName:  "User2",
		Password:  "password123",
	}

	regBody2, _ := json.Marshal(registerReq2)
	regReq2 := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(regBody2))
	regReq2.Header.Set("Content-Type", "application/json")
	_, err := suite.app.Test(regReq2)
	if err != nil {
		return
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/", nil)
	req.Header.Set("Cookie", "session_token="+sessionToken)

	resp, err := suite.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response dto.UsersListResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return
	}
	assert.Len(suite.T(), response.Users, 2)
}

func extractSessionTokenFromCookie(cookieHeader string) string {
	parts := strings.Split(cookieHeader, ";")
	for _, part := range parts {
		if strings.HasPrefix(strings.TrimSpace(part), "session_token=") {
			return strings.TrimPrefix(strings.TrimSpace(part), "session_token=")
		}
	}
	return ""
}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}
