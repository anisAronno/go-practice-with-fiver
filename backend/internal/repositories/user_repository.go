package repositories

import (
	"gofiver/internal/models"

	"gorm.io/gorm"
)

// UserRepository handles user database operations (like Laravel Repository pattern)
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindAll returns all users with pagination
func (r *UserRepository) FindAll(offset, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	r.db.Model(&models.User{}).Count(&total)
	err := r.db.Offset(offset).Limit(limit).Find(&users).Error

	return users, total, err
}

// Update updates a user
func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete soft deletes a user
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// ExistsByEmail checks if email exists
func (r *UserRepository) ExistsByEmail(email string, excludeID ...uint) bool {
	var count int64
	query := r.db.Model(&models.User{}).Where("email = ?", email)
	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}
	query.Count(&count)
	return count > 0
}
