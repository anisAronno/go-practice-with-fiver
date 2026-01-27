package services

import (
	"errors"
	"gofiver/internal/dto"
	"gofiver/internal/models"
	"gofiver/internal/repositories"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetAll(pagination *dto.PaginationQuery) ([]models.User, int64, error) {
	return s.userRepo.FindAll(pagination.GetOffset(), pagination.GetPerPage())
}

func (s *UserService) GetByID(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *UserService) Update(id uint, req *dto.UpdateUserRequest) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

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

func (s *UserService) UpdateRole(id uint, role string) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	user.Role = role

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Delete(id uint) error {
	_, err := s.userRepo.FindByID(id)
	if err != nil {
		return errors.New("user not found")
	}
	return s.userRepo.Delete(id)
}

func (s *UserService) GetAllDeleted(pagination *dto.PaginationQuery) ([]models.User, int64, error) {
	return s.userRepo.FindAllDeleted(pagination.GetOffset(), pagination.GetPerPage())
}

func (s *UserService) Restore(id uint) error {
	_, err := s.userRepo.FindDeletedByID(id)
	if err != nil {
		return errors.New("deleted user not found")
	}
	return s.userRepo.Restore(id)
}

func (s *UserService) ForceDelete(id uint) error {
	_, err := s.userRepo.FindDeletedByID(id)
	if err != nil {
		return errors.New("deleted user not found")
	}
	return s.userRepo.ForceDelete(id)
}
