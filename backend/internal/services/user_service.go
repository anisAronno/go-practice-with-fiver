package services

import (
	"errors"
	"gofiver/internal/dto"
	"gofiver/internal/models"
	"gofiver/internal/repositories"
)

// UserService handles user business logic
type UserService struct {
	userRepo *repositories.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// GetAll returns paginated users
func (s *UserService) GetAll(pagination *dto.PaginationQuery) ([]models.User, int64, error) {
	return s.userRepo.FindAll(pagination.GetOffset(), pagination.GetPerPage())
}

// GetByID returns a user by ID
func (s *UserService) GetByID(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

// Update updates a user
func (s *UserService) Update(id uint, req *dto.UpdateUserRequest) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Check email uniqueness if changing
	if req.Email != "" && req.Email != user.Email {
		if s.userRepo.ExistsByEmail(req.Email, id) {
			return nil, errors.New("email already exists")
		}
		user.Email = req.Email
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Delete deletes a user
func (s *UserService) Delete(id uint) error {
	_, err := s.userRepo.FindByID(id)
	if err != nil {
		return errors.New("user not found")
	}
	return s.userRepo.Delete(id)
}
