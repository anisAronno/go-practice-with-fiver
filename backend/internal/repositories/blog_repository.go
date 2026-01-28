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
	err := r.db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name, email")
	}).First(&blog, id).Error
	if err != nil {
		return nil, err
	}
	return &blog, nil
}

func (r *BlogRepository) FindAll(offset, limit int, search string) ([]models.Blog, error) {
	var blogs []models.Blog
	query := r.db.Select("blogs.id, blogs.title, blogs.image, blogs.user_id, blogs.created_at").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).
		Order("blogs.id DESC").
		Offset(offset).
		Limit(limit)

	if search != "" {
		query = query.Where("blogs.title LIKE ?", search+"%")
	}

	err := query.Find(&blogs).Error
	return blogs, err
}

func (r *BlogRepository) GetCount(search string) (int64, error) {
	var count int64
	query := r.db.Model(&models.Blog{})
	if search != "" {
		query = query.Where("title LIKE ?", search+"%")
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *BlogRepository) FindByUserID(userID uint, offset, limit int) ([]models.Blog, error) {
	var blogs []models.Blog
	err := r.db.Select("id, title, image, user_id, created_at").
		Where("user_id = ?", userID).
		Order("id DESC").
		Offset(offset).
		Limit(limit).
		Find(&blogs).Error
	return blogs, err
}

func (r *BlogRepository) GetUserBlogCount(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Blog{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
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

func (r *BlogRepository) FindAllDeleted(offset, limit int) ([]models.Blog, error) {
	var blogs []models.Blog
	err := r.db.Unscoped().
		Select("blogs.id, blogs.title, blogs.image, blogs.user_id, blogs.created_at, blogs.deleted_at").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).
		Where("blogs.deleted_at IS NOT NULL").
		Order("blogs.deleted_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&blogs).Error
	return blogs, err
}

func (r *BlogRepository) GetDeletedCount() (int64, error) {
	var count int64
	err := r.db.Unscoped().Model(&models.Blog{}).Where("deleted_at IS NOT NULL").Count(&count).Error
	return count, err
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
