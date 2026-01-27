package repositories

import (
	"gofiver/internal/models"

	"gorm.io/gorm"
)

// BlogRepository handles blog database operations
type BlogRepository struct {
	db *gorm.DB
}

// NewBlogRepository creates a new blog repository
func NewBlogRepository(db *gorm.DB) *BlogRepository {
	return &BlogRepository{db: db}
}

// Create creates a new blog
func (r *BlogRepository) Create(blog *models.Blog) error {
	return r.db.Create(blog).Error
}

// FindByID finds a blog by ID with user
func (r *BlogRepository) FindByID(id uint) (*models.Blog, error) {
	var blog models.Blog
	err := r.db.Preload("User").First(&blog, id).Error
	if err != nil {
		return nil, err
	}
	return &blog, nil
}

// FindAll returns all blogs with pagination
func (r *BlogRepository) FindAll(offset, limit int) ([]models.Blog, int64, error) {
	var blogs []models.Blog
	var total int64

	r.db.Model(&models.Blog{}).Count(&total)
	err := r.db.Preload("User").Offset(offset).Limit(limit).Order("created_at DESC").Find(&blogs).Error

	return blogs, total, err
}

// FindByUserID returns blogs by user ID
func (r *BlogRepository) FindByUserID(userID uint, offset, limit int) ([]models.Blog, int64, error) {
	var blogs []models.Blog
	var total int64

	r.db.Model(&models.Blog{}).Where("user_id = ?", userID).Count(&total)
	err := r.db.Where("user_id = ?", userID).Offset(offset).Limit(limit).Order("created_at DESC").Find(&blogs).Error

	return blogs, total, err
}

// Update updates a blog
func (r *BlogRepository) Update(blog *models.Blog) error {
	return r.db.Save(blog).Error
}

// Delete deletes a blog
func (r *BlogRepository) Delete(id uint) error {
	return r.db.Delete(&models.Blog{}, id).Error
}

// IsOwner checks if user owns the blog
func (r *BlogRepository) IsOwner(blogID, userID uint) bool {
	var count int64
	r.db.Model(&models.Blog{}).Where("id = ? AND user_id = ?", blogID, userID).Count(&count)
	return count > 0
}
