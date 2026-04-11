package dto

import (
	"time"

	"github.com/google/uuid"
)

type WishListRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventDate   time.Time `json:"event_date"`
}

type WishListUpdateRequest struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	EventDate   *time.Time `json:"event_date,omitempty"`
}

type WishListResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventDate   time.Time `json:"event_date"`
	ShareToken  string    `json:"share_token"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type WishListListResponse struct {
	WishLists []WishListResponse `json:"wishlists"`
}
