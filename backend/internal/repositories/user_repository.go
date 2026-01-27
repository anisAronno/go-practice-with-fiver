package repositories

import (
	"gofiver/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindAll(offset, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	r.db.Model(&models.User{}).Count(&total)
	err := r.db.Offset(offset).Limit(limit).Find(&users).Error

	return users, total, err
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *UserRepository) ExistsByEmail(email string, excludeID ...uint) bool {
	var count int64
	query := r.db.Model(&models.User{}).Where("email = ?", email)
	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}
	query.Count(&count)
	return count > 0
}

func (r *UserRepository) FindAllDeleted(offset, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	r.db.Unscoped().Model(&models.User{}).Where("deleted_at IS NOT NULL").Count(&total)
	err := r.db.Unscoped().Where("deleted_at IS NOT NULL").Offset(offset).Limit(limit).Order("deleted_at DESC").Find(&users).Error

	return users, total, err
}

func (r *UserRepository) FindDeletedByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Unscoped().Where("id = ? AND deleted_at IS NOT NULL", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Restore(id uint) error {
	return r.db.Unscoped().Model(&models.User{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *UserRepository) ForceDelete(id uint) error {
	return r.db.Unscoped().Delete(&models.User{}, id).Error
}
