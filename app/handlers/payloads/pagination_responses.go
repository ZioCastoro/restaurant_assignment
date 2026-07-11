package payloads

import (
	"github.com/ZioCastoro/restaurant_assignment/app/entities"
	"github.com/ZioCastoro/restaurant_assignment/utils"
)

type PaginationResponse struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize,omitempty"`
}

func NewPaginationResponse(entity entities.Pagination) PaginationResponse {
	return PaginationResponse{
		PageSize: utils.SafeDereference(entity.PageSize),
		Page:     utils.SafeDereference(entity.Page, 1),
	}
}
