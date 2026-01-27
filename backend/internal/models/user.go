package models

import (
	"golang.org/x/crypto/bcrypt"
)

// User represents the users table (like Laravel's User model)
type User struct {
	BaseModel
	Name     string `json:"name" gorm:"size:255;not null"`
	Email    string `json:"email" gorm:"size:255;uniqueIndex;not null"`
	Password string `json:"-" gorm:"size:255;not null"`
	Blogs    []Blog `json:"blogs,omitempty" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name
func (User) TableName() string {
	return "users"
}

// HashPassword hashes the user password
func (u *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword verifies the password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// ToResponse returns user without sensitive data
func (u *User) ToResponse() map[string]interface{} {
	return map[string]interface{}{
		"id":         u.ID,
		"name":       u.Name,
		"email":      u.Email,
		"created_at": u.CreatedAt,
		"updated_at": u.UpdatedAt,
	}
}
