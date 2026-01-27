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

type AuthService struct {
	userRepo *repositories.UserRepository
	config   *config.Config
}

func NewAuthService(userRepo *repositories.UserRepository, config *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		config:   config,
	}
}

func (s *AuthService) Register(req *dto.RegisterRequest) (*models.User, string, error) {

	if s.userRepo.ExistsByEmail(req.Email) {
		return nil, "", errors.New("email already exists")
	}


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


	token, err := s.generateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

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

func (s *AuthService) generateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(s.config.JWT.Expiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.Secret))
}

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
