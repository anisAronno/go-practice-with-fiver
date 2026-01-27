package dto

// RegisterRequest represents registration data (like Laravel FormRequest)
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=255"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginRequest represents login data
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UpdateUserRequest represents user update data
type UpdateUserRequest struct {
	Name  string `json:"name" validate:"omitempty,min=2,max=255"`
	Email string `json:"email" validate:"omitempty,email"`
}
