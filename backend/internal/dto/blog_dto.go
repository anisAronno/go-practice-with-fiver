package dto

type CreateBlogRequest struct {
	Title   string `json:"title" validate:"required,min=3,max=255"`
	Content string `json:"content" validate:"required,min=10"`
}

type UpdateBlogRequest struct {
	Title   string `json:"title" validate:"omitempty,min=3,max=255"`
	Content string `json:"content" validate:"omitempty,min=10"`
}

type BlogSearchQuery struct {
	Page    int    `query:"page"`
	PerPage int    `query:"per_page"`
	Search  string `query:"search"`
}

func (q *BlogSearchQuery) GetPage() int {
	if q.Page < 1 {
		return 1
	}
	return q.Page
}

func (q *BlogSearchQuery) GetPerPage() int {
	if q.PerPage < 1 {
		return 20
	}
	if q.PerPage > 100 {
		return 100
	}
	return q.PerPage
}

func (q *BlogSearchQuery) GetOffset() int {
	return (q.GetPage() - 1) * q.GetPerPage()
}
