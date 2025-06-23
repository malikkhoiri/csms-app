package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/malikkhoiri/csms/internal/config"
	"github.com/malikkhoiri/csms/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  domain.UserRepository
	jwtConfig *config.JWTConfig
}

func NewAuthService(userRepo domain.UserRepository, jwtConfig *config.JWTConfig) domain.AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtConfig: jwtConfig,
	}
}

func (s *AuthService) Login(ctx context.Context, request *domain.AuthRequest) (*domain.AuthResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, request.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check if user is active
	if user.Status != "active" {
		return nil, errors.New("user account is not active")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	// Clear password from response
	user.Password = ""

	return &domain.AuthResponse{
		AccessToken: token,
		User:        *user,
	}, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (*domain.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.jwtConfig.Secret), nil
	})

	if err != nil {
		return nil, errors.New("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Extract user ID from claims
		userID, ok := claims["user_id"].(float64)
		if !ok {
			return nil, errors.New("invalid token claims")
		}

		user, err := s.userRepo.GetByID(ctx, uint(userID))
		if err != nil {
			return nil, errors.New("user not found")
		}

		if user.Status != "active" {
			return nil, errors.New("user account is not active")
		}

		user.Password = ""
		return user, nil
	}

	return nil, errors.New("invalid token")
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	// For now, just return success
	return nil
}

func (s *AuthService) generateToken(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(s.jwtConfig.Expiration).Unix(),
		"iat":     time.Now().Unix(),
		"iss":     s.jwtConfig.Issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtConfig.Secret))
}

func (s *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
