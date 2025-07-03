package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"linkedin-watcher/internal/models"
	"linkedin-watcher/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthController_Register_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Create a real auth service for testing
	authService := services.NewAuthService(nil, "test-secret")
	authController := NewAuthController(authService)

	auth := router.Group("/auth")
	{
		auth.POST("/register", authController.Register)
	}

	// Invalid JSON
	body := []byte(`{"email": "invalid-json"`)
	request := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthController_Login_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Create a real auth service for testing
	authService := services.NewAuthService(nil, "test-secret")
	authController := NewAuthController(authService)

	auth := router.Group("/auth")
	{
		auth.POST("/login", authController.Login)
	}

	// Invalid JSON
	body := []byte(`{"email": "invalid-json"`)
	request := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthController_RefreshToken_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Create a real auth service for testing
	authService := services.NewAuthService(nil, "test-secret")
	authController := NewAuthController(authService)

	auth := router.Group("/auth")
	{
		auth.POST("/refresh", authController.RefreshToken)
	}

	// Invalid JSON
	body := []byte(`{"refresh_token": "invalid-json"`)
	request := httptest.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAuthController_Logout_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Create a real auth service for testing
	authService := services.NewAuthService(nil, "test-secret")
	authController := NewAuthController(authService)

	auth := router.Group("/auth")
	{
		auth.POST("/logout", authController.Logout)
	}

	request := httptest.NewRequest("POST", "/auth/logout", nil)
	request.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Logged out successfully", response["message"])
}

func TestNewAuthController(t *testing.T) {
	authService := &services.AuthService{}
	controller := NewAuthController(authService)
	assert.NotNil(t, controller)
	assert.Equal(t, authService, controller.authService)
}

func TestAuthController_UserInfoStructureNeverContainsID(t *testing.T) {
	// Test that our UserInfo model structure doesn't contain ID field
	userInfo := models.UserInfo{
		Email:    "test@example.com",
		Name:     "Test User",
		AuthType: "password",
	}

	// Marshal to JSON and verify structure
	jsonData, err := json.Marshal(userInfo)
	assert.NoError(t, err)

	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonData, &jsonMap)
	assert.NoError(t, err)

	// Verify that "id" field is not present
	_, hasID := jsonMap["id"]
	assert.False(t, hasID, "UserInfo should not contain 'id' field in JSON response")

	// Verify expected fields are present
	assert.Contains(t, jsonMap, "email")
	assert.Contains(t, jsonMap, "name")
	assert.Contains(t, jsonMap, "auth_type")
}

func TestAuthController_AuthResponseStructureNeverContainsID(t *testing.T) {
	// Test that our AuthResponse model structure doesn't contain ID field in User
	authResponse := models.AuthResponse{
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
		User: models.UserInfo{
			Email:    "test@example.com",
			Name:     "Test User",
			AuthType: "password",
		},
	}

	// Marshal to JSON and verify structure
	jsonData, err := json.Marshal(authResponse)
	assert.NoError(t, err)

	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonData, &jsonMap)
	assert.NoError(t, err)

	// Verify that user object doesn't contain "id" field
	userObj, hasUser := jsonMap["user"].(map[string]interface{})
	assert.True(t, hasUser, "AuthResponse should contain 'user' field")

	_, hasID := userObj["id"]
	assert.False(t, hasID, "User object in AuthResponse should not contain 'id' field")

	// Verify expected fields are present
	assert.Contains(t, userObj, "email")
	assert.Contains(t, userObj, "name")
	assert.Contains(t, userObj, "auth_type")
}

func TestAuthController_RequestModelsStructure(t *testing.T) {
	// Test that our request models have the correct structure
	registerData := models.UserRegistration{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
	}

	loginData := models.UserLogin{
		Email:    "test@example.com",
		Password: "password123",
	}

	passwordChangeData := models.PasswordChangeRequest{
		CurrentPassword: "oldpassword",
		NewPassword:     "newpassword123",
	}

	refreshTokenData := models.RefreshTokenRequest{
		RefreshToken: "refresh_token_here",
	}

	// Verify none of these contain ID fields
	registerJSON, _ := json.Marshal(registerData)
	loginJSON, _ := json.Marshal(loginData)
	passwordChangeJSON, _ := json.Marshal(passwordChangeData)
	refreshTokenJSON, _ := json.Marshal(refreshTokenData)

	var registerMap, loginMap, passwordChangeMap, refreshTokenMap map[string]interface{}

	json.Unmarshal(registerJSON, &registerMap)
	json.Unmarshal(loginJSON, &loginMap)
	json.Unmarshal(passwordChangeJSON, &passwordChangeMap)
	json.Unmarshal(refreshTokenJSON, &refreshTokenMap)

	// Verify no ID fields in any request models
	_, hasID := registerMap["id"]
	assert.False(t, hasID, "UserRegistration should not contain 'id' field")

	_, hasID = loginMap["id"]
	assert.False(t, hasID, "UserLogin should not contain 'id' field")

	_, hasID = passwordChangeMap["id"]
	assert.False(t, hasID, "PasswordChangeRequest should not contain 'id' field")

	_, hasID = refreshTokenMap["id"]
	assert.False(t, hasID, "RefreshTokenRequest should not contain 'id' field")
}
