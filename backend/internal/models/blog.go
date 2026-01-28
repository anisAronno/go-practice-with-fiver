package models

type Blog struct {
	BaseModel
	Title   string  `json:"title" gorm:"size:255;not null;index:idx_blog_title"`
	Content string  `json:"content" gorm:"type:text;not null"`
	Image   *string `json:"image" gorm:"size:500"`
	UserID  uint    `json:"user_id" gorm:"not null;index:idx_blog_user_id;index:idx_blog_user_created,priority:1"`
	User    User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (Blog) TableName() string {
	return "blogs"
}

type BlogListItem struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Image     string `json:"image"`
	UserID    uint   `json:"user_id"`
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
	CreatedAt string `json:"created_at"`
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

func (b *Blog) ToListResponse() map[string]interface{} {
	response := map[string]interface{}{
		"id":         b.ID,
		"title":      b.Title,
		"image":      b.Image,
		"user_id":    b.UserID,
		"created_at": b.CreatedAt,
	}

	if b.User.ID != 0 {
		response["user"] = map[string]interface{}{
			"id":   b.User.ID,
			"name": b.User.Name,
		}
	}

	return response
}
