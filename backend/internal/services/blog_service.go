package services

import (
	"errors"
	"gofiver/internal/dto"
	"gofiver/internal/models"
	"gofiver/internal/repositories"
)

// BlogService handles blog business logic
type BlogService struct {
	blogRepo *repositories.BlogRepository
}

// NewBlogService creates a new blog service
func NewBlogService(blogRepo *repositories.BlogRepository) *BlogService {
	return &BlogService{blogRepo: blogRepo}
}

// Create creates a new blog
func (s *BlogService) Create(userID uint, req *dto.CreateBlogRequest) (*models.Blog, error) {
	blog := &models.Blog{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID,
	}

	if err := s.blogRepo.Create(blog); err != nil {
		return nil, err
	}

	// Reload with user
	return s.blogRepo.FindByID(blog.ID)
}

// GetAll returns paginated blogs
func (s *BlogService) GetAll(pagination *dto.PaginationQuery) ([]models.Blog, int64, error) {
	return s.blogRepo.FindAll(pagination.GetOffset(), pagination.GetPerPage())
}

// GetByID returns a blog by ID
func (s *BlogService) GetByID(id uint) (*models.Blog, error) {
	return s.blogRepo.FindByID(id)
}

// GetByUserID returns blogs by user
func (s *BlogService) GetByUserID(userID uint, pagination *dto.PaginationQuery) ([]models.Blog, int64, error) {
	return s.blogRepo.FindByUserID(userID, pagination.GetOffset(), pagination.GetPerPage())
}

// Update updates a blog (only owner can update)
func (s *BlogService) Update(id, userID uint, req *dto.UpdateBlogRequest) (*models.Blog, error) {
	blog, err := s.blogRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("blog not found")
	}

	// Check ownership
	if blog.UserID != userID {
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

// Delete deletes a blog (only owner can delete)
func (s *BlogService) Delete(id, userID uint) error {
	blog, err := s.blogRepo.FindByID(id)
	if err != nil {
		return errors.New("blog not found")
	}

	// Check ownership
	if blog.UserID != userID {
		return errors.New("unauthorized")
	}

	return s.blogRepo.Delete(id)
}
