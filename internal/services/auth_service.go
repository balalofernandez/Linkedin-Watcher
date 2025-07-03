package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"linkedin-watcher/db"
	"linkedin-watcher/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	queries   *db.Queries
	jwtSecret []byte
}

func NewAuthService(queries *db.Queries, jwtSecret string) *AuthService {
	return &AuthService{
		queries:   queries,
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *AuthService) Register(ctx context.Context, req models.UserRegistration) (*models.AuthResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := s.queries.CreateUserWithPassword(ctx, db.CreateUserWithPasswordParams{
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: pgtype.Text{String: string(hashedPassword), Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	userID := uuid.UUID(user.ID.Bytes).String()
	accessToken, refreshToken, err := s.generateTokens(userID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: models.UserInfo{
			Email:    user.Email,
			Name:     user.Name,
			AuthType: user.AuthType,
		},
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req models.UserLogin) (*models.AuthResponse, error) {
	user, err := s.queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	if !user.PasswordHash.Valid {
		return nil, errors.New("account uses OAuth authentication")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	userID := uuid.UUID(user.ID.Bytes).String()
	accessToken, refreshToken, err := s.generateTokens(userID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: models.UserInfo{
			Email:    user.Email,
			Name:     user.Name,
			AuthType: user.AuthType,
		},
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*models.AuthResponse, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(*models.JWTClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userUUID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, errors.New("invalid user ID in token")
	}

	pgUUID := pgtype.UUID{Bytes: userUUID, Valid: true}
	user, err := s.queries.GetUserByID(ctx, pgUUID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	userID := uuid.UUID(user.ID.Bytes).String()
	accessToken, newRefreshToken, err := s.generateTokens(userID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User: models.UserInfo{
			Email:    user.Email,
			Name:     user.Name,
			AuthType: user.AuthType,
		},
	}, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*models.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*models.JWTClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID string, currentPassword, newPassword string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	pgUUID := pgtype.UUID{Bytes: userUUID, Valid: true}
	user, err := s.queries.GetUserByID(ctx, pgUUID)
	if err != nil {
		return errors.New("user not found")
	}

	if !user.PasswordHash.Valid {
		return errors.New("account uses OAuth authentication")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash.String), []byte(currentPassword)); err != nil {
		return errors.New("invalid current password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	err = s.queries.UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
		ID:           user.ID,
		PasswordHash: pgtype.Text{String: string(hashedPassword), Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

func (s *AuthService) generateTokens(userID string, email string) (string, string, error) {
	now := time.Now()
	accessExp := now.Add(15 * time.Minute)
	refreshExp := now.Add(7 * 24 * time.Hour)

	accessClaims := &models.JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExp),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	refreshClaims := &models.JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExp),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	accessTokenString, err := accessToken.SignedString(s.jwtSecret)
	if err != nil {
		return "", "", err
	}

	refreshTokenString, err := refreshToken.SignedString(s.jwtSecret)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func (s *AuthService) generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
