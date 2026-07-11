package entities

type Pagination struct {
	Page     *int
	PageSize *int
}

func (p *Pagination) GetLimit() *int {
	return p.PageSize
}

func (p *Pagination) GetOffset() int {
	if p.PageSize == nil || p.Page == nil {
		return 0
	}

	return (*p.Page - 1) * *p.PageSize
}
