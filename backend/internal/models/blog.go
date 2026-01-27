package models

// Blog represents the blogs table
type Blog struct {
	BaseModel
	Title   string `json:"title" gorm:"size:255;not null"`
	Content string `json:"content" gorm:"type:text;not null"`
	UserID  uint   `json:"user_id" gorm:"not null;index"`
	User    User   `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name
func (Blog) TableName() string {
	return "blogs"
}

// ToResponse returns blog with user info
func (b *Blog) ToResponse() map[string]interface{} {
	response := map[string]interface{}{
		"id":         b.ID,
		"title":      b.Title,
		"content":    b.Content,
		"user_id":    b.UserID,
		"created_at": b.CreatedAt,
		"updated_at": b.UpdatedAt,
	}

	if b.User.ID != 0 {
		response["user"] = b.User.ToResponse()
	}

	return response
}
