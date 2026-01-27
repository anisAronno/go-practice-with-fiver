package services

import (
	"errors"
	"gofiver/internal/config"
	"gofiver/internal/dto"
	"gofiver/internal/models"
	"gofiver/internal/repositories"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AuthService handles authentication logic
type AuthService struct {
	userRepo *repositories.UserRepository
	config   *config.Config
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo *repositories.UserRepository, config *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		config:   config,
	}
}

// Register creates a new user account
func (s *AuthService) Register(req *dto.RegisterRequest) (*models.User, string, error) {
	// Check if email exists
	if s.userRepo.ExistsByEmail(req.Email) {
		return nil, "", errors.New("email already exists")
	}

	// Create user
	user := &models.User{
		Name:  req.Name,
		Email: req.Email,
	}

	if err := user.HashPassword(req.Password); err != nil {
		return nil, "", err
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, "", err
	}

	// Generate token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// Login authenticates a user
func (s *AuthService) Login(req *dto.LoginRequest) (*models.User, string, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if !user.CheckPassword(req.Password) {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// generateToken creates a JWT token
func (s *AuthService) generateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(s.config.JWT.Expiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.Secret))
}

// ValidateToken validates and parses a JWT token
func (s *AuthService) ValidateToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return &claims, nil
}
