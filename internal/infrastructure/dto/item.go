package dto

import (
	"time"

	"github.com/google/uuid"
)

type WishListItemRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	ProductURL  *string `json:"product_url,omitempty"`
	Priority    int     `json:"priority"`
}

type WishListItemUpdateRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	ProductURL  *string `json:"product_url,omitempty"`
	Priority    *int    `json:"priority,omitempty"`
}

type WishListItemResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description,omitempty"`
	ProductURL  *string   `json:"product_url,omitempty"`
	Priority    int       `json:"priority"`
	IsReserved  bool      `json:"is_reserved"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type WishListItemListResponse struct {
	Items []WishListItemResponse `json:"items"`
}
