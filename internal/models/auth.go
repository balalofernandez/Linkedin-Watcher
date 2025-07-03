package models

import (
	"github.com/golang-jwt/jwt/v5"
)

// AuthType represents the type of authentication used
type AuthType string

const (
	AuthTypePassword AuthType = "password"
	AuthTypeGoogle   AuthType = "google"
	AuthTypeBoth     AuthType = "both"
)

// UserRegistration represents the data needed to register a new user
type UserRegistration struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

// UserLogin represents the data needed to login a user
type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// JWTClaims represents the claims in a JWT token
type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// AuthResponse represents the response from authentication operations
type AuthResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token,omitempty"`
	User         UserInfo `json:"user"`
}

// UserInfo represents user information returned in responses
type UserInfo struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	AuthType string `json:"auth_type"`
}

// PasswordChangeRequest represents the data needed to change password
type PasswordChangeRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

// RefreshTokenRequest represents the data needed to refresh a token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
