package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/AtlasOpx/devprep/internal/config"
	"github.com/AtlasOpx/devprep/internal/database"
	"github.com/AtlasOpx/devprep/internal/dto"
	"github.com/AtlasOpx/devprep/internal/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthFlowTestSuite struct {
	suite.Suite
	app *fiber.App
	db  *database.DB
	cfg *config.Config
}

func (suite *AuthFlowTestSuite) SetupSuite() {
	suite.cfg = &config.Config{
		DatabaseURL: "postgres://postgres:password@localhost:5432/devprep_test?sslmode=disable",
		ServerPort:  "3001",
	}

	var err error
	suite.db, err = database.Connect(suite.cfg)
	if err != nil {
		suite.T().Skipf("Database connection failed: %v", err)
	}

	suite.app = fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "DevPrep",
		AppName:       "Dev Prep app v1.0.1",
		ReadTimeout:   30 * time.Second,
		WriteTimeout:  30 * time.Second,
		IdleTimeout:   120 * time.Second,
	})

	suite.app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	routes.SetupRoutes(suite.app, suite.db, suite.cfg)
}

func (suite *AuthFlowTestSuite) TearDownSuite() {
	if suite.db != nil {
		suite.db.Close()
	}
}

func (suite *AuthFlowTestSuite) SetupTest() {
	suite.cleanupDatabase()
}

func (suite *AuthFlowTestSuite) TearDownTest() {
	suite.cleanupDatabase()
}

// TODO(): ПОхУй
func (suite *AuthFlowTestSuite) cleanupDatabase() {
	if suite.db != nil {
		_, _ = suite.db.DB.Exec("DELETE FROM sessions")
		_, _ = suite.db.DB.Exec("DELETE FROM users")
	}
}

func (suite *AuthFlowTestSuite) TestCompleteUserRegistrationAndLoginFlow() {
	registerReq := dto.RegisterRequest{
		Email:     "testuser@example.com",
		Username:  "testuser123",
		FirstName: "Test",
		LastName:  "User",
		Password:  "securepassword123",
	}

	reqBody, _ := json.Marshal(registerReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.app.Test(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var registerResponse dto.RegisterResponse
	err = json.NewDecoder(resp.Body).Decode(&registerResponse)
	if err != nil {
		return
	}
	assert.Equal(suite.T(), "User created successfully", registerResponse.Message)
	assert.NotEmpty(suite.T(), registerResponse.UserID)

	loginReq := dto.LoginRequest{
		Email:    registerReq.Email,
		Password: registerReq.Password,
	}

	loginBody, _ := json.Marshal(loginReq)
	loginRequest := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(loginBody))
	loginRequest.Header.Set("Content-Type", "application/json")

	loginResp, err := suite.app.Test(loginRequest)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, loginResp.StatusCode)

	var loginResponse dto.LoginResponse
	err = json.NewDecoder(loginResp.Body).Decode(&loginResponse)
	if err != nil {
		return
	}
	assert.Equal(suite.T(), "Login successful", loginResponse.Message)
	assert.Equal(suite.T(), registerReq.Email, loginResponse.User.Email)
	assert.Equal(suite.T(), registerReq.Username, loginResponse.User.Username)

	cookies := loginResp.Header.Get("Set-Cookie")
	assert.Contains(suite.T(), cookies, "session_token")

	sessionToken := extractSessionToken(cookies)
	assert.NotEmpty(suite.T(), sessionToken)

	profileReq := httptest.NewRequest(http.MethodGet, "/api/v1/users/profile", nil)
	profileReq.Header.Set("Cookie", fmt.Sprintf("session_token=%s", sessionToken))

	profileResp, err := suite.app.Test(profileReq)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, profileResp.StatusCode)

	var profileResponse dto.UserResponse
	err = json.NewDecoder(profileResp.Body).Decode(&profileResponse)
	if err != nil {
		return
	}
	assert.Equal(suite.T(), registerReq.Email, profileResponse.Email)
	assert.Equal(suite.T(), registerReq.Username, profileResponse.Username)
	assert.Equal(suite.T(), registerReq.FirstName, profileResponse.FirstName)
	assert.Equal(suite.T(), registerReq.LastName, profileResponse.LastName)

	updateReq := dto.UpdateProfileRequest{
		FirstName: "Updated",
		LastName:  "Name",
		Username:  "updateduser",
	}

	updateBody, _ := json.Marshal(updateReq)
	updateRequest := httptest.NewRequest(http.MethodPut, "/api/v1/users/profile", bytes.NewBuffer(updateBody))
	updateRequest.Header.Set("Content-Type", "application/json")
	updateRequest.Header.Set("Cookie", fmt.Sprintf("session_token=%s", sessionToken))

	updateResp, err := suite.app.Test(updateRequest)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, updateResp.StatusCode)

	var updateResponse dto.UpdateProfileResponse
	err = json.NewDecoder(updateResp.Body).Decode(&updateResponse)
	if err != nil {
		return
	}
	assert.Equal(suite.T(), "Profile updated successfully", updateResponse.Message)

	verifyReq := httptest.NewRequest(http.MethodGet, "/api/v1/users/profile", nil)
	verifyReq.Header.Set("Cookie", fmt.Sprintf("session_token=%s", sessionToken))

	verifyResp, err := suite.app.Test(verifyReq)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, verifyResp.StatusCode)

	var verifyResponse dto.UserResponse
	err = json.NewDecoder(verifyResp.Body).Decode(&verifyResponse)
	if err != nil {
		return
	}
	assert.Equal(suite.T(), updateReq.FirstName, verifyResponse.FirstName)
	assert.Equal(suite.T(), updateReq.LastName, verifyResponse.LastName)
	assert.Equal(suite.T(), updateReq.Username, verifyResponse.Username)

	logoutReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	logoutReq.Header.Set("Cookie", fmt.Sprintf("session_token=%s", sessionToken))

	logoutResp, err := suite.app.Test(logoutReq)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, logoutResp.StatusCode)

	var logoutResponse dto.LogoutResponse
	err = json.NewDecoder(logoutResp.Body).Decode(&logoutResponse)
	if err != nil {
		return
	}
	assert.Equal(suite.T(), "Logout successful", logoutResponse.Message)

	unauthorizedReq := httptest.NewRequest(http.MethodGet, "/api/v1/users/profile", nil)
	unauthorizedReq.Header.Set("Cookie", fmt.Sprintf("session_token=%s", sessionToken))

	unauthorizedResp, err := suite.app.Test(unauthorizedReq)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusUnauthorized, unauthorizedResp.StatusCode)
}

func (suite *AuthFlowTestSuite) TestMultipleUsersRegistrationAndListUsers() {
	users := []dto.RegisterRequest{
		{
			Email:     "user1@example.com",
			Username:  "user1",
			FirstName: "User",
			LastName:  "One",
			Password:  "password123",
		},
		{
			Email:     "user2@example.com",
			Username:  "user2",
			FirstName: "User",
			LastName:  "Two",
			Password:  "password123",
		},
		{
			Email:     "user3@example.com",
			Username:  "user3",
			FirstName: "User",
			LastName:  "Three",
			Password:  "password123",
		},
	}

	for _, user := range users {
		reqBody, _ := json.Marshal(user)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := suite.app.Test(req)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)
	}

	loginReq := dto.LoginRequest{
		Email:    users[0].Email,
		Password: users[0].Password,
	}

	loginBody, _ := json.Marshal(loginReq)
	loginRequest := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(loginBody))
	loginRequest.Header.Set("Content-Type", "application/json")

	loginResp, err := suite.app.Test(loginRequest)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, loginResp.StatusCode)

	sessionToken := extractSessionToken(loginResp.Header.Get("Set-Cookie"))

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/users/", nil)
	listReq.Header.Set("Cookie", fmt.Sprintf("session_token=%s", sessionToken))

	listResp, err := suite.app.Test(listReq)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, listResp.StatusCode)

	var usersResponse dto.UsersListResponse
	err = json.NewDecoder(listResp.Body).Decode(&usersResponse)
	if err != nil {
		return
	}
	assert.Len(suite.T(), usersResponse.Users, 3)

	emails := make(map[string]bool)
	for _, user := range usersResponse.Users {
		emails[user.Email] = true
	}

	for _, expectedUser := range users {
		assert.True(suite.T(), emails[expectedUser.Email], "Expected user %s not found in response", expectedUser.Email)
	}
}

func (suite *AuthFlowTestSuite) TestInvalidRegistrationAttempts() {
	testCases := []struct {
		name     string
		request  dto.RegisterRequest
		expected int
	}{
		{
			name: "Invalid email format",
			request: dto.RegisterRequest{
				Email:     "invalid-email",
				Username:  "validuser",
				FirstName: "Test",
				LastName:  "User",
				Password:  "password123",
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "Empty username",
			request: dto.RegisterRequest{
				Email:     "test@example.com",
				Username:  "",
				FirstName: "Test",
				LastName:  "User",
				Password:  "password123",
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "Short password",
			request: dto.RegisterRequest{
				Email:     "test@example.com",
				Username:  "validuser",
				FirstName: "Test",
				LastName:  "User",
				Password:  "123",
			},
			expected: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tc.request)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			resp, err := suite.app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, resp.StatusCode)
		})
	}
}

func (suite *AuthFlowTestSuite) TestDuplicateEmailRegistration() {
	registerReq := dto.RegisterRequest{
		Email:     "duplicate@example.com",
		Username:  "user1",
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

	registerReq.Username = "user2"
	reqBody2, _ := json.Marshal(registerReq)
	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(reqBody2))
	req2.Header.Set("Content-Type", "application/json")

	resp2, err := suite.app.Test(req2)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusBadRequest, resp2.StatusCode)
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

func TestAuthFlowTestSuite(t *testing.T) {
	suite.Run(t, new(AuthFlowTestSuite))
}
