package repositories

import (
	"gofiver/internal/models"

	"gorm.io/gorm"
)

type BlogRepository struct {
	db *gorm.DB
}

func NewBlogRepository(db *gorm.DB) *BlogRepository {
	return &BlogRepository{db: db}
}

func (r *BlogRepository) Create(blog *models.Blog) error {
	return r.db.Create(blog).Error
}

func (r *BlogRepository) FindByID(id uint) (*models.Blog, error) {
	var blog models.Blog
	err := r.db.Preload("User").First(&blog, id).Error
	if err != nil {
		return nil, err
	}
	return &blog, nil
}

func (r *BlogRepository) FindAll(offset, limit int) ([]models.Blog, int64, error) {
	var blogs []models.Blog
	var total int64

	r.db.Model(&models.Blog{}).Count(&total)

	err := r.db.
		Select("blogs.id, blogs.title, blogs.image, blogs.user_id, blogs.created_at").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).
		Offset(offset).
		Limit(limit).
		Order("blogs.id DESC").
		Find(&blogs).Error

	return blogs, total, err
}

func (r *BlogRepository) FindByUserID(userID uint, offset, limit int) ([]models.Blog, int64, error) {
	var blogs []models.Blog
	var total int64

	r.db.Model(&models.Blog{}).Where("user_id = ?", userID).Count(&total)

	err := r.db.
		Select("id, title, image, user_id, created_at").
		Where("user_id = ?", userID).
		Offset(offset).
		Limit(limit).
		Order("id DESC").
		Find(&blogs).Error

	return blogs, total, err
}

func (r *BlogRepository) Update(blog *models.Blog) error {
	return r.db.Save(blog).Error
}

func (r *BlogRepository) Delete(id uint) error {
	return r.db.Delete(&models.Blog{}, id).Error
}

func (r *BlogRepository) IsOwner(blogID, userID uint) bool {
	var count int64
	r.db.Model(&models.Blog{}).Where("id = ? AND user_id = ?", blogID, userID).Count(&count)
	return count > 0
}

func (r *BlogRepository) FindAllDeleted(offset, limit int) ([]models.Blog, int64, error) {
	var blogs []models.Blog
	var total int64

	r.db.Unscoped().Model(&models.Blog{}).Where("deleted_at IS NOT NULL").Count(&total)

	err := r.db.Unscoped().
		Select("blogs.id, blogs.title, blogs.image, blogs.user_id, blogs.created_at, blogs.deleted_at").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).
		Where("blogs.deleted_at IS NOT NULL").
		Offset(offset).
		Limit(limit).
		Order("blogs.deleted_at DESC").
		Find(&blogs).Error

	return blogs, total, err
}

func (r *BlogRepository) FindDeletedByID(id uint) (*models.Blog, error) {
	var blog models.Blog
	err := r.db.Unscoped().Preload("User").Where("id = ? AND deleted_at IS NOT NULL", id).First(&blog).Error
	if err != nil {
		return nil, err
	}
	return &blog, nil
}

func (r *BlogRepository) Restore(id uint) error {
	return r.db.Unscoped().Model(&models.Blog{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *BlogRepository) ForceDelete(id uint) error {
	return r.db.Unscoped().Delete(&models.Blog{}, id).Error
}

func (r *BlogRepository) UpdateImage(id uint, imagePath string) error {
	return r.db.Model(&models.Blog{}).Where("id = ?", id).Update("image", imagePath).Error
}
