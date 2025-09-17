package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AtlasOpx/devprep/internal/config"
	"github.com/AtlasOpx/devprep/internal/database"
	"github.com/AtlasOpx/devprep/internal/dto"
	"github.com/AtlasOpx/devprep/internal/handlers"
	"github.com/AtlasOpx/devprep/internal/repository"
	"github.com/AtlasOpx/devprep/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthHandlerTestSuite struct {
	suite.Suite
	app         *fiber.App
	db          *database.DB
	cfg         *config.Config
	userRepo    *repository.UserRepository
	authRepo    *repository.AuthRepository
	authService *service.AuthService
	authHandler *handlers.AuthHandler
}

func (suite *AuthHandlerTestSuite) SetupSuite() {
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
	suite.authService = service.NewAuthService(suite.userRepo, suite.authRepo)
	suite.authHandler = handlers.NewAuthHandler(suite.authService, suite.authRepo, suite.cfg)

	suite.app = fiber.New()
	suite.setupRoutes()
}

func (suite *AuthHandlerTestSuite) TearDownSuite() {
	if suite.db != nil {
		suite.db.Close()
	}
}

func (suite *AuthHandlerTestSuite) SetupTest() {
	suite.cleanupDatabase()
}

func (suite *AuthHandlerTestSuite) TearDownTest() {
	suite.cleanupDatabase()
}

func (suite *AuthHandlerTestSuite) cleanupDatabase() {
	if suite.db != nil {
		if _, err := suite.db.DB.Exec("DELETE FROM sessions"); err != nil {
			suite.T().Logf("Failed to delete sessions: %v", err)
		}

		if _, err := suite.db.DB.Exec("DELETE FROM users"); err != nil {
			suite.T().Logf("Failed to delete users: %v", err)
		}
	}
}

func (suite *AuthHandlerTestSuite) setupRoutes() {
	api := suite.app.Group("/api/v1")
	auth := api.Group("/auth")

	auth.Post("/register", suite.authHandler.Register)
	auth.Post("/login", suite.authHandler.Login)
	auth.Post("/logout", suite.authHandler.Logout)
}

func (suite *AuthHandlerTestSuite) TestRegister_Success() {
	registerReq := dto.RegisterRequest{
		Email:     "test@example.com",
		Username:  "testuser",
		FirstName: "Test",
		LastName:  "User",
		Password:  "password123",
	}

	reqBody, _ := json.Marshal(registerReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var response dto.RegisterResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Errorf("Failed to decode response: %v", err)
	}
	assert.Equal(suite.T(), "User created successfully", response.Message)
	assert.NotEmpty(suite.T(), response.UserID)
}

func (suite *AuthHandlerTestSuite) TestRegister_InvalidInput() {
	registerReq := dto.RegisterRequest{
		Email:    "invalid-email",
		Username: "",
		Password: "123",
	}

	reqBody, _ := json.Marshal(registerReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

func (suite *AuthHandlerTestSuite) TestRegister_DuplicateEmail() {
	registerReq := dto.RegisterRequest{
		Email:     "test@example.com",
		Username:  "testuser1",
		FirstName: "Test",
		LastName:  "User",
		Password:  "password123",
	}

	reqBody, _ := json.Marshal(registerReq)

	req1 := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(reqBody))
	req1.Header.Set("Content-Type", "application/json")

	resp1, err := suite.app.Test(req1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusCreated, resp1.StatusCode)

	registerReq.Username = "testuser2"
	reqBody2, _ := json.Marshal(registerReq)

	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(reqBody2))
	req2.Header.Set("Content-Type", "application/json")

	resp2, err := suite.app.Test(req2)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusBadRequest, resp2.StatusCode)
}

func (suite *AuthHandlerTestSuite) TestLogin_Success() {
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
	_, err := suite.app.Test(regReq)
	if err != nil {
		return
	}

	loginReq := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	reqBody, _ := json.Marshal(loginReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response dto.LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Errorf("%v", err)
	}
	assert.Equal(suite.T(), "Login successful", response.Message)
	assert.NotEmpty(suite.T(), response.User.ID)

	cookies := resp.Header.Get("Set-Cookie")
	assert.Contains(suite.T(), cookies, "session_token")
}

func (suite *AuthHandlerTestSuite) TestLogin_InvalidCredentials() {
	loginReq := dto.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "wrongpassword",
	}

	reqBody, _ := json.Marshal(loginReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
}

func (suite *AuthHandlerTestSuite) TestLogin_WrongPassword() {
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
	_, err := suite.app.Test(regReq)
	if err != nil {
		return
	}

	loginReq := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	reqBody, _ := json.Marshal(loginReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
}

func (suite *AuthHandlerTestSuite) TestLogout_Success() {
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
	_, err := suite.app.Test(regReq)
	if err != nil {
		return
	}

	loginReq := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	loginBody, _ := json.Marshal(loginReq)
	loginReqHttp := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(loginBody))
	loginReqHttp.Header.Set("Content-Type", "application/json")
	loginResp, _ := suite.app.Test(loginReqHttp)

	sessionToken := extractSessionToken(loginResp.Header.Get("Set-Cookie"))

	logoutReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	logoutReq.Header.Set("Cookie", fmt.Sprintf("session_token=%s", sessionToken))

	resp, err := suite.app.Test(logoutReq)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var response dto.LogoutResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return
	}
	assert.Equal(suite.T(), "Logout successful", response.Message)
}

func (suite *AuthHandlerTestSuite) TestLogout_NoSession() {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)

	resp, err := suite.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusBadRequest, resp.StatusCode)
}

func extractSessionToken(cookieHeader string) string {
	parts := strings.Split(cookieHeader, ";")
	for _, part := range parts {
		if strings.HasPrefix(strings.TrimSpace(part), "session_token=") {
			return strings.TrimPrefix(strings.TrimSpace(part), "session_token=")
		}
	}
	return ""
}

func TestAuthHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(AuthHandlerTestSuite))
}
