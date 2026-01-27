package models

type Blog struct {
	BaseModel
	Title   string  `json:"title" gorm:"size:255;not null"`
	Content string  `json:"content" gorm:"type:text;not null"`
	Image   *string `json:"image" gorm:"size:500"`
	UserID  uint    `json:"user_id" gorm:"not null;index"`
	User    User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (Blog) TableName() string {
	return "blogs"
}

func (b *Blog) ToResponse() map[string]interface{} {
	response := map[string]interface{}{
		"id":         b.ID,
		"title":      b.Title,
		"content":    b.Content,
		"image":      b.Image,
		"user_id":    b.UserID,
		"created_at": b.CreatedAt,
		"updated_at": b.UpdatedAt,
	}

	if b.DeletedAt.Valid {
		response["deleted_at"] = b.DeletedAt.Time
	}

	if b.User.ID != 0 {
		response["user"] = b.User.ToResponse()
	}

	return response
}
