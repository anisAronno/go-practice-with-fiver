package dto

type PaginationQuery struct {
	Page    int `query:"page" validate:"omitempty,min=1"`
	PerPage int `query:"per_page" validate:"omitempty,min=1,max=100"`
}

func (p *PaginationQuery) GetPage() int {
	if p.Page < 1 {
		return 1
	}
	return p.Page
}

func (p *PaginationQuery) GetPerPage() int {
	if p.PerPage < 1 {
		return 15
	}
	return p.PerPage
}

func (p *PaginationQuery) GetOffset() int {
	return (p.GetPage() - 1) * p.GetPerPage()
}
