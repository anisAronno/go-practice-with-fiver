package dto

type CreateBlogRequest struct {
	Title   string `json:"title" validate:"required,min=3,max=255"`
	Content string `json:"content" validate:"required,min=10"`
}

type UpdateBlogRequest struct {
	Title   string `json:"title" validate:"omitempty,min=3,max=255"`
	Content string `json:"content" validate:"omitempty,min=10"`
}
