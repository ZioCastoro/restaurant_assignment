package payloads

import "github.com/ZioCastoro/restaurant_assignment/app/entities"

const (
	DefaultPage     = 1
	DefaultPageSize = 25
	MaxPageSize     = 100
)

type PaginationRequest struct {
	Page     int `query:"page,omitempty"`
	PageSize int `query:"page_size,omitempty"`
}

func (p *PaginationRequest) MapToEntity() entities.Pagination {
	pagination := entities.Pagination{
		Page:     new(DefaultPage),
		PageSize: new(DefaultPageSize),
	}

	if p.Page > 0 {
		pagination.Page = &p.Page
	}

	if p.PageSize > 0 && p.PageSize <= MaxPageSize {
		pagination.PageSize = &p.PageSize
	}

	return pagination
}
