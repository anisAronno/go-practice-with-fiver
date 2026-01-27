package services

import (
	"errors"
	"gofiver/internal/dto"
	"gofiver/internal/models"
	"gofiver/internal/repositories"
)

type BlogService struct {
	blogRepo *repositories.BlogRepository
}

func NewBlogService(blogRepo *repositories.BlogRepository) *BlogService {
	return &BlogService{blogRepo: blogRepo}
}

func (s *BlogService) Create(userID uint, req *dto.CreateBlogRequest) (*models.Blog, error) {
	blog := &models.Blog{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID,
	}

	if err := s.blogRepo.Create(blog); err != nil {
		return nil, err
	}

	return s.blogRepo.FindByID(blog.ID)
}

func (s *BlogService) GetAll(pagination *dto.PaginationQuery) ([]models.Blog, int64, error) {
	return s.blogRepo.FindAll(pagination.GetOffset(), pagination.GetPerPage())
}

func (s *BlogService) GetByID(id uint) (*models.Blog, error) {
	return s.blogRepo.FindByID(id)
}

func (s *BlogService) GetByUserID(userID uint, pagination *dto.PaginationQuery) ([]models.Blog, int64, error) {
	return s.blogRepo.FindByUserID(userID, pagination.GetOffset(), pagination.GetPerPage())
}

func (s *BlogService) Update(id, userID uint, req *dto.UpdateBlogRequest, isAdmin bool) (*models.Blog, error) {
	blog, err := s.blogRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("blog not found")
	}

	if blog.UserID != userID && !isAdmin {
		return nil, errors.New("unauthorized")
	}

	if req.Title != "" {
		blog.Title = req.Title
	}
	if req.Content != "" {
		blog.Content = req.Content
	}

	if err := s.blogRepo.Update(blog); err != nil {
		return nil, err
	}

	return blog, nil
}

func (s *BlogService) Delete(id, userID uint, isAdmin bool) error {
	blog, err := s.blogRepo.FindByID(id)
	if err != nil {
		return errors.New("blog not found")
	}

	if blog.UserID != userID && !isAdmin {
		return errors.New("unauthorized")
	}

	return s.blogRepo.Delete(id)
}

func (s *BlogService) GetAllDeleted(pagination *dto.PaginationQuery) ([]models.Blog, int64, error) {
	return s.blogRepo.FindAllDeleted(pagination.GetOffset(), pagination.GetPerPage())
}

func (s *BlogService) Restore(id uint) error {
	_, err := s.blogRepo.FindDeletedByID(id)
	if err != nil {
		return errors.New("deleted blog not found")
	}
	return s.blogRepo.Restore(id)
}

func (s *BlogService) ForceDelete(id uint) error {
	_, err := s.blogRepo.FindDeletedByID(id)
	if err != nil {
		return errors.New("deleted blog not found")
	}
	return s.blogRepo.ForceDelete(id)
}

func (s *BlogService) UpdateImage(id uint, imagePath string) error {
	blog, err := s.blogRepo.FindByID(id)
	if err != nil {
		return errors.New("blog not found")
	}
	blog.Image = &imagePath
	return s.blogRepo.Update(blog)
}
