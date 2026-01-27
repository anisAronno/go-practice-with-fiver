package models

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	BaseModel
	Name     string `json:"name" gorm:"size:255;not null"`
	Email    string `json:"email" gorm:"size:255;uniqueIndex;not null"`
	Password string `json:"-" gorm:"size:255;not null"`
	Role     string `json:"role" gorm:"size:50;default:author"`
	Blogs    []Blog `json:"blogs,omitempty" gorm:"foreignKey:UserID"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

func (u *User) ToResponse() map[string]interface{} {
	return map[string]interface{}{
		"id":         u.ID,
		"name":       u.Name,
		"email":      u.Email,
		"role":       u.Role,
		"created_at": u.CreatedAt,
		"updated_at": u.UpdatedAt,
	}
}
