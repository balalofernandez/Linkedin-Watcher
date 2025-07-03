package controllers

import (
	"linkedin-watcher/internal/models"
	"linkedin-watcher/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthController handles authentication-related HTTP requests
type AuthController struct {
	authService *services.AuthService
}

// NewAuthController creates a new AuthController with injected dependencies
func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.UserRegistration true "User registration data"
// @Success 201 {object} models.AuthResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /auth/register [post]
func (ac *AuthController) Register(c *gin.Context) {
	var req models.UserRegistration
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	response, err := ac.authService.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// @Summary Login user
// @Description Login user with email and password (supports both JSON body and HTTP Basic Auth)
// @Tags auth
// @Accept json
// @Produce json
// @Security BasicAuth
// @Param user body models.UserLogin false "User login data (optional if using Basic Auth)"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/login [post]
func (ac *AuthController) Login(c *gin.Context) {
	var req models.UserLogin

	// Try to get credentials from Basic Auth header first
	email, password, hasBasicAuth := c.Request.BasicAuth()

	if hasBasicAuth {
		// Use Basic Auth credentials
		req = models.UserLogin{
			Email:    email,
			Password: password,
		}
	} else {
		// Try to bind JSON body
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid request data. Provide either JSON body or HTTP Basic Auth",
				"details": err.Error(),
			})
			return
		}
	}

	response, err := ac.authService.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh_token body models.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/refresh [post]
func (ac *AuthController) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	response, err := ac.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Change password
// @Description Change user password
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param password_change body models.PasswordChangeRequest true "Password change data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/change-password [post]
func (ac *AuthController) ChangePassword(c *gin.Context) {
	var req models.PasswordChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	userID := c.MustGet("userID").(string)

	err := ac.authService.ChangePassword(c.Request.Context(), userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password changed successfully",
	})
}

// @Summary Logout user
// @Description Logout user (client should discard tokens)
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /auth/logout [post]
func (ac *AuthController) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}
