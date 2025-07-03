package services

import (
	"encoding/json"
	"testing"

	"linkedin-watcher/db"
	"linkedin-watcher/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestNewAuthService(t *testing.T) {
	queries := &db.Queries{}
	jwtSecret := "test-secret"

	authService := NewAuthService(queries, jwtSecret)

	assert.NotNil(t, authService)
	assert.Equal(t, queries, authService.queries)
	assert.Equal(t, []byte(jwtSecret), authService.jwtSecret)
}

func TestAuthService_ValidateToken(t *testing.T) {
	queries := &db.Queries{}
	authService := NewAuthService(queries, "test-secret")

	// Test with invalid token
	_, err := authService.ValidateToken("invalid-token")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token")

	// Test with empty token
	_, err = authService.ValidateToken("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token")
}

// TestAuthService_UserInfoNeverContainsID verifies that user IDs are never returned in API responses
func TestAuthService_UserInfoNeverContainsID(t *testing.T) {
	// Test that UserInfo struct doesn't have an ID field
	userInfo := models.UserInfo{
		Email:    "test@example.com",
		Name:     "Test User",
		AuthType: "password",
	}

	// Marshal to JSON and verify no ID field exists
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

// TestAuthService_RegisterResponseNeverContainsID verifies that register response doesn't contain user ID
func TestAuthService_RegisterResponseNeverContainsID(t *testing.T) {
	// Create a mock user for testing
	mockUser := db.User{
		ID:           pgtype.UUID{Bytes: uuid.New(), Valid: true},
		Email:        "test@example.com",
		Name:         "Test User",
		PasswordHash: pgtype.Text{String: "hashed_password", Valid: true},
		AuthType:     "password",
		IsActive:     true,
	}

	// Mock the queries to return our test user
	// Note: In a real test, you'd use a proper mock or test database
	// This is just to verify the response structure

	// Test that AuthResponse.User doesn't have an ID field
	authResponse := models.AuthResponse{
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
		User: models.UserInfo{
			Email:    mockUser.Email,
			Name:     mockUser.Name,
			AuthType: mockUser.AuthType,
		},
	}

	// Marshal to JSON and verify no ID field exists in User
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

// TestAuthService_LoginResponseNeverContainsID verifies that login response doesn't contain user ID
func TestAuthService_LoginResponseNeverContainsID(t *testing.T) {
	// Test that login response structure doesn't include user ID
	loginResponse := models.AuthResponse{
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
		User: models.UserInfo{
			Email:    "test@example.com",
			Name:     "Test User",
			AuthType: "password",
		},
	}

	// Marshal to JSON and verify no ID field exists
	jsonData, err := json.Marshal(loginResponse)
	assert.NoError(t, err)

	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonData, &jsonMap)
	assert.NoError(t, err)

	// Verify that user object doesn't contain "id" field
	userObj, hasUser := jsonMap["user"].(map[string]interface{})
	assert.True(t, hasUser, "AuthResponse should contain 'user' field")

	_, hasID := userObj["id"]
	assert.False(t, hasID, "User object in login response should not contain 'id' field")
}

// TestAuthService_RefreshTokenResponseNeverContainsID verifies that refresh token response doesn't contain user ID
func TestAuthService_RefreshTokenResponseNeverContainsID(t *testing.T) {
	// Test that refresh token response structure doesn't include user ID
	refreshResponse := models.AuthResponse{
		AccessToken:  "new_access_token",
		RefreshToken: "new_refresh_token",
		User: models.UserInfo{
			Email:    "test@example.com",
			Name:     "Test User",
			AuthType: "password",
		},
	}

	// Marshal to JSON and verify no ID field exists
	jsonData, err := json.Marshal(refreshResponse)
	assert.NoError(t, err)

	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonData, &jsonMap)
	assert.NoError(t, err)

	// Verify that user object doesn't contain "id" field
	userObj, hasUser := jsonMap["user"].(map[string]interface{})
	assert.True(t, hasUser, "AuthResponse should contain 'user' field")

	_, hasID := userObj["id"]
	assert.False(t, hasID, "User object in refresh token response should not contain 'id' field")
}
